# Reproducible Alpine test-image fixtures for the billable acceptance tests
# (BRIEFING.md §8.3). Fetches a pinned Alpine nocloud raw image, bakes a
# throwaway SSH public key into /root/.ssh/authorized_keys, enables sshd + DHCP,
# and recompresses to .raw.xz. Fully hermetic: the fetch is fixed-output and the
# bake runs offline via libguestfs (no root, no network).
#
# The disk edits use guestfish, so no privileged loop-mount is required and the
# build works in the Nix sandbox.
{
  pkgs,
  # "x86_64" | "aarch64" — selects the Alpine artifact and asserts the arch.
  arch,
  # The throwaway public key baked into the image (path to the .pub file).
  authorizedKey,
}:
let
  # Pinned Alpine nocloud generic-cloud images. Bump version + hash together;
  # obtain a hash by setting it to pkgs.lib.fakeHash and reading the mismatch.
  alpineVersion = "3.20.3";
  alpineBranch = "v3.20";

  images = {
    "x86_64" = {
      url = "https://dl-cdn.alpinelinux.org/alpine/${alpineBranch}/releases/cloud/nocloud_alpine-${alpineVersion}-x86_64-bios-cloudinit-r0.qcow2";
      hash = "sha256-CUfDkfW/TzBZtJmgwy8WKniG3BAsHn7O6W1l7YDpmbI=";
    };
    "aarch64" = {
      url = "https://dl-cdn.alpinelinux.org/alpine/${alpineBranch}/releases/cloud/nocloud_alpine-${alpineVersion}-aarch64-uefi-cloudinit-r0.qcow2";
      hash = "sha256-vu2hVaeIRSuH+9UqUa5NrNN8m2sLUQVrlYis3gSKT7E=";
    };
  };

  img = images.${arch};

  src = pkgs.fetchurl {
    inherit (img) url hash;
  };

  keyContents = builtins.readFile authorizedKey;
in
pkgs.stdenv.mkDerivation {
  pname = "hcloudimage-test-image-${arch}";
  version = alpineVersion;

  nativeBuildInputs = [
    # guestfs-tools provides virt-customize / virt-inspector; libguestfs-with-appliance
    # supplies the QEMU appliance they drive to edit the image offline. The appliance
    # needs /dev/kvm at build time, so this derivation builds on a KVM-capable machine
    # (the acceptance runner) — see README.
    pkgs.guestfs-tools
    pkgs.libguestfs-with-appliance
    pkgs.qemu-utils
    pkgs.xz
  ];

  # Point guestfs tools at the bundled appliance instead of trying to build one.
  LIBGUESTFS_PATH = "${pkgs.libguestfs-with-appliance}/lib/guestfs";

  dontUnpack = true;

  # Deterministic compression.
  SOURCE_DATE_EPOCH = "1";

  buildPhase = ''
    runHook preBuild

    export HOME=$TMPDIR
    export LIBGUESTFS_BACKEND=direct
    export LIBGUESTFS_CACHEDIR=$TMPDIR

    # Work on a writable qcow2, convert to raw, then customise offline.
    cp ${src} disk.qcow2
    chmod +w disk.qcow2
    qemu-img convert -f qcow2 -O raw disk.qcow2 image.raw

    # Bake the throwaway key + enable sshd and DHCP, all offline via the
    # libguestfs appliance (virt-customize auto-detects the root partition).
    printf '%s' ${pkgs.lib.escapeShellArg keyContents} > authorized_keys
    printf 'auto eth0\niface eth0 inet dhcp\n' > interfaces

    # --no-logfile and a fixed timestamp reduce in-image nondeterminism; disk
    # images are not bit-reproducible in general (fs metadata/journals), but the
    # inputs are pinned and the customize step is fully offline/hermetic.
    virt-customize -a image.raw \
      --no-logfile \
      --mkdir /root/.ssh \
      --upload authorized_keys:/root/.ssh/authorized_keys \
      --chmod 0700:/root/.ssh \
      --chmod 0600:/root/.ssh/authorized_keys \
      --upload interfaces:/etc/network/interfaces \
      --run-command 'rc-update add sshd default || true' \
      --run-command 'rc-update add networking default || true'

    # Assert the guest architecture matches (no silent x86 fallback for arm).
    # Write to a file first: piping virt-inspector into grep -q makes grep close
    # the pipe on first match, which crashes the inspector with a broken pipe.
    virt-inspector -a image.raw > inspector.xml
    if ! grep -q '<arch>${arch}</arch>' inspector.xml; then
      echo "fixture arch is not ${arch}" >&2
      exit 1
    fi

    xz -T0 -6 image.raw

    runHook postBuild
  '';

  installPhase = ''
    runHook preInstall
    mkdir -p $out
    cp image.raw.xz $out/test-image-${arch}.raw.xz
    runHook postInstall
  '';

  meta = {
    description = "Reproducible Alpine ${arch} acceptance fixture (.raw.xz) with a baked throwaway SSH key";
    license = pkgs.lib.licenses.mit;
  };
}
