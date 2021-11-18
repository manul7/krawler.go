1. I use "github.com/urfave/cli/v2" instead of "github.com/spf13/cobra" due to higher complexity of Cobra utility (it's lib and generator)
2. CLI framework allows easy addition of flags if needed. For example, I can pass flag to clean existing destination directory without asking.
3. I use "filepath" instead of "path" to achieve portability, it's recommeded.
4. I use "github.com/PuerkitoBio/goquery" to parse HTML content and extract links, standard net.html library looks too low level.