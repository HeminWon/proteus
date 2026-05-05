class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.13"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.13/proteus-v0.1.13-darwin-arm64.tar.gz"
      sha256 "b713b884f5f2d333262f6e3f6193dca88669bfe16438abfc1b3d7e4e867d6252"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.13/proteus-v0.1.13-linux-arm64.tar.gz"
      sha256 "df38ab910bf531f4bfd8fca8859e952e7327a3428d18bf2feb166af939a44d8b"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.13/proteus-v0.1.13-linux-x64.tar.gz"
      sha256 "44a6d723860d0310dde6011fb5534f955f49a80e7c7c148cbe692006d6a6a925"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
