class Purtypics < Formula
  desc "A fast, modern static photo gallery generator"
  homepage "https://github.com/chrissnell/purtypics"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/chrissnell/purtypics/releases/download/v#{version}/purtypics_#{version}_macOS_arm64.tar.gz"
      sha256 "SHA256_ARM64"
    else
      url "https://github.com/chrissnell/purtypics/releases/download/v#{version}/purtypics_#{version}_macOS_x86_64.tar.gz"
      sha256 "SHA256_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/chrissnell/purtypics/releases/download/v#{version}/purtypics_#{version}_linux_arm64.tar.gz"
      sha256 "SHA256_LINUX_ARM64"
    else
      url "https://github.com/chrissnell/purtypics/releases/download/v#{version}/purtypics_#{version}_linux_x86_64.tar.gz"
      sha256 "SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "purtypics"
    (share/"purtypics/themes").install Dir["themes/*"] if Dir.exist?("themes")
  end

  test do
    assert_match "purtypics", shell_output("#{bin}/purtypics --version 2>&1", 0)
  end
end
