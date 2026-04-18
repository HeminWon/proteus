class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.0.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.0.0/proteus-v0.0.0-darwin-arm64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.0.0/proteus-v0.0.0-darwin-x64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.0.0/proteus-v0.0.0-linux-arm64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.0.0/proteus-v0.0.0-linux-x64.tar.gz"
      sha256 "0000000000000000000000000000000000000000000000000000000000000000"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
