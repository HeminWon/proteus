class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.9"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.9/proteus-v0.1.9-darwin-arm64.tar.gz"
      sha256 "66c290426750638d2ccdc54ff6c3e2158cdebbb75680988fbb1e034e1890d4e7"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.9/proteus-v0.1.9-linux-arm64.tar.gz"
      sha256 "76a2df84449b2a9a24c46ca02b7e55c8037e215f61cfa34ca954eba2dcb46e81"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.9/proteus-v0.1.9-linux-x64.tar.gz"
      sha256 "d37d41296434cb602885544bc12689cbb5c5a921ca1c57560975f73ead1e35f9"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
