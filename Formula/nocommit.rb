class Nocommit < Formula
  desc "Prevent accidental commits with forbidden phrases"
  homepage "https://github.com/waleedbaroudi/nocommit"
  version "0.1.0" #TODO: inject?

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/waleedbaroudi/nocommit/releases/download/v0.1.0/nocommit_darwin_arm64.tar.gz"
      sha256 "ba13cff3e26108095103b1d3658b0a742f820bda533917843e96ba25245f6444"
    else
      url "https://github.com/waleedbaroudi/nocommit/releases/download/v0.1.0/nocommit_darwin_amd64.tar.gz"
      sha256 "503441aa92dfa45f9af3917ac5ea99d68369bff04597c5f45c33764a86fa6917"
    end
  end

  def install
    bin.install "nocommit"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/nocommit --version")
  end
end
