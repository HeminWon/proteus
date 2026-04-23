class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.7"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.7/proteus-v0.1.7-darwin-arm64.tar.gz"
      sha256 "cc42b87e328e6ec6d1632cf3e6c59b765d4f7f3db0161d9f8701c9376ee61a8f"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.7/proteus-v0.1.7-linux-arm64.tar.gz"
      sha256 "c6641374a352a1f08a8726639123bf4fcb468e98658d1d0c96b0b9f0335f57a4"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.7/proteus-v0.1.7-linux-x64.tar.gz"
      sha256 "19aa2ea42e81a18be152e5f433128168b45547857e5698bd74a4fc080e1336e3"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
