class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.3"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.3/proteus-v0.1.3-darwin-arm64.tar.gz"
      sha256 "76ff077d08a3b566569e81d8bef2bbbbe2200902a4b5f9cf77ec4805cd457083"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.3/proteus-v0.1.3-linux-arm64.tar.gz"
      sha256 "72af6649dd7b48f4c65489dcd15b9d2260fb8ce1326f2272598e29b8df4f2a71"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.3/proteus-v0.1.3-linux-x64.tar.gz"
      sha256 "d1b50817f8af61dae01c93fe5b1375235aa8b73613b2217936bc48287f9cb774"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
