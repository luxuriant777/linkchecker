# Link Checker
Link Checker is a command-line tool written in Go for checking broken links on a website.

## Installation
To use Link Checker, you need to have Go installed on your machine. If you haven't installed Go, please follow the
official Go installation guide: [https://go.dev/doc/install](https://go.dev/doc/install)

Once you have Go installed, you can follow these steps to install Link Checker:

1. Clone the repository:
   ```shell
   git clone https://github.com/luxuriant777/linkchecker
   ```
2. Navigate to the project directory:
    ```shell
    cd linkchecker
    ```
3. Build the project:
    ```shell
    go build ./cmd/linkchecker
    ```
   This will create an executable file named `linkchecker` in the current directory.


4. Run the Link Checker:
   - Linux:
    ```shell
    ./linkchecker <url>
    ```
   - Windows:
    ```shell
    linkchecker.exe <url>
    ```
    Replace `<url>` with the URL of the website you want to check for broken links, for example:
    ```shell
    ./linkchecker https://example.com
    ```

## Usage
To use Link Checker, follow these steps:

1. Open your terminal or command prompt.

2. Navigate to the directory where you have the linkchecker executable.

3. Run the Link Checker:
   - Linux:
    ```shell
    ./linkchecker <url>
    ```
   - Windows:
    ```shell
    linkchecker.exe <url>
    ```
   Replace `<url>` with the URL of the website you want to check for broken links, for example:
    ```shell
    ./linkchecker https://example.com
    ```

4. The Link Checker will scan the webpage, extract links, and check the status of each link. 
The results will be displayed in the terminal.

    Example output:

    ```shell
    Link: https://example.com/ Status: OK
    Link: https://example.com/page1 Status: Broken
    Link: https://example.com/page2 Status: OK
    ```

    That's it! You can now use Link Checker to check for broken links on websites.

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please create an issue 
or submit a pull request on the `dev` branch.

## License

This project is licensed under the
[GNU General Public License (GPL) version 3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).
