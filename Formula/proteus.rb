class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.5"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.5/proteus-v0.1.5-darwin-arm64.tar.gz"
      sha256 "4c59cb500450424b39ae1bcb9d16a13ac4359f20bae784d5d24040a5303bcdc9"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.5/proteus-v0.1.5-linux-arm64.tar.gz"
      sha256 "cddbe95351a4e8a19d6533e5c746b405327c5fc385ef84f6db43e227d564f8ce"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.5/proteus-v0.1.5-linux-x64.tar.gz"
      sha256 "15bcad4a154bb32d873586084228c1a63ad8a0de786285a806738ee8bb01add3"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
