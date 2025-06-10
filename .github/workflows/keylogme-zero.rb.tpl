# .github/workflows/myproject.rb.tpl

class Myproject < Formula
  desc "A keylogger to record your keypresses stats"
  homepage "https://www.keylogme.com"
  license "MIT"

  # The version will be updated dynamically from the GitHub Release tag
  version "{{TAG_NAME_NO_V}}"

  # Define URLs and SHA256s conditionally for each OS/Architecture
  on_macos do
    on_intel do
      url "{{URL_DARWIN_AMD64}}"
      sha256 "{{SHA256_DARWIN_AMD64}}"
    end
    on_arm do
      url "{{URL_DARWIN_ARM64}}"
      sha256 "{{SHA256_DARWIN_ARM64}}"
    end
  end

  on_linux do
    on_intel do
      url "{{URL_LINUX_AMD64}}"
      sha256 "{{SHA256_LINUX_AMD64}}"
    end
    on_arm do
      url "{{URL_LINUX_ARM64}}"
      sha256 "{{SHA256_LINUX_ARM64}}"
    end
  end

  def install
    # Assuming your downloaded tarball/zip contains the binary directly at the root
    # or in a predictable subdirectory within the extracted archive.
    # If your binary is named differently for each platform (e.g., myproject-mac, myproject-linux),
    # you might need conditional `bin.install` here as well, but usually it's named the same.
    bin.install keylogme-zero
  end

  test do
    # Simple test to ensure the binary runs and exits successfully
    # IMPORTANT: Replace "myproject" with your actual binary name
    system "#{bin}/keylogme-zero", "--version" # Or any simple command that returns 0
  end
end
