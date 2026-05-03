class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.8"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.8/proteus-v0.1.8-darwin-arm64.tar.gz"
      sha256 "4772aec0b8a0cde77f103da073da9ee9522c5b80342177b852038ff8b7a6993a"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.8/proteus-v0.1.8-linux-arm64.tar.gz"
      sha256 "eabf6e923afa791556252562b6a8269dbab4d0d31873b6bcd46cae9116f33b3a"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.8/proteus-v0.1.8-linux-x64.tar.gz"
      sha256 "b19ce52d88cc74ec8a5b153c450f81222b3c261c48f4f385fd1aee5f5e8ed34f"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
