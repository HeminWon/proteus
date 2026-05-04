class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.11"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.11/proteus-v0.1.11-darwin-arm64.tar.gz"
      sha256 "fc8a10c30c9d0f9003a2f1c3de286aaa9d198b7425c6c00a28852f256e083333"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.11/proteus-v0.1.11-linux-arm64.tar.gz"
      sha256 "4b2e712f35d233b8dd99ccae0e1efc099555ec179bc8f8fa9bd44501af42ca10"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.11/proteus-v0.1.11-linux-x64.tar.gz"
      sha256 "f2bbf1093b4d6ec91e55ebc933d89e27bb97d8e0c8d59667b28d3088d68130b2"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
