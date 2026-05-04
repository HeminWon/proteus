class Proteus < Formula
  desc "Switch providers, models, and config profiles for AI developer tools"
  homepage "https://github.com/HeminWon/proteus"
  version "0.1.10"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.10/proteus-v0.1.10-darwin-arm64.tar.gz"
      sha256 "3d3442cc3e7604ae5f1a705e54255f9abe2fbfd471bea785d5dbb85503715c63"
    else
      odie "No darwin-x64 release asset is available for this version."
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.10/proteus-v0.1.10-linux-arm64.tar.gz"
      sha256 "30c0eb0c3e15375fd2232261b9866d6e42905c52d8f4dc9da51abdf844f40975"
    else
      url "https://github.com/HeminWon/proteus/releases/download/v0.1.10/proteus-v0.1.10-linux-x64.tar.gz"
      sha256 "67cee6406df701ed8bb1f6023f00adfc3cb6266c46561d2f4cc10bfdc45d710c"
    end
  end

  def install
    bin.install "proteus"
  end

  test do
    assert_match "proteus", shell_output("#{bin}/proteus --help")
  end
end
