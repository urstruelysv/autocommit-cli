class AutocommitCli < Formula
  desc "An AI-powered commit message generator"
  homepage "https://github.com/urstruelysv/autocommit-cli"
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "install", "./cmd/autocommit-cli"
  end

  test do
    # Basic test to ensure the CLI runs
    assert_match "Welcome to autocommit-cli!", shell_output("#{bin}/autocommit-cli --help")
  end
end
