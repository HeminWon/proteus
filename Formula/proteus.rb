class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.12"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.12/proteus-v0.1.12-darwin-arm64.tar.gz"
      sha256 "9c698586e38a5741a852946036619914ba8b6ce876cfbd08ba2a93018dbb2d35"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.12/proteus-v0.1.12-linux-arm64.tar.gz"
      sha256 "991f9a0bf4e330ebe15992d6fbe52fcbdf5f5aedae7be0ca396d260350e4774d"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.12/proteus-v0.1.12-linux-x64.tar.gz"
      sha256 "3aeffc1221bde46e76c9e94f32cea26f1e5f4512394b709124678a086373a781"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
