# Link Checker
Link Checker is a powerful command-line tool, written in Go, designed to improve the health and accessibility
of your website. It provides robust features not only for identifying broken links, but also for tracking the
status codes returned by all URLs on your site. The results are conveniently sorted and saved in separate files
according to their corresponding HTTP status codes.
![Screenshot_2](https://github.com/luxuriant777/linkchecker/assets/20545475/e014c523-d7a5-4c73-a52a-a387d3e9abf1)

In addition to checking links and logging response statuses, Link Checker can also serve as an invaluable tool
for creating a comprehensive list of all the URLS available on the website. By effectively traversing and 
cataloguing every reachable URL on a website, it enables you to create comprehensive sitemaps. This functionality
can significantly aid in improving the SEO of your site, making your content more discoverable and navigable.
![Screenshot_3](https://github.com/luxuriant777/linkchecker/assets/20545475/53bbb500-7ee3-411b-8f41-58c9cdc2edf5)

## Features
- Broken Link Checking: Identifies and logs all broken links (HTTP status 404) on your website. Broken links
can negatively impact your website's user experience and SEO. All broken links can be found in `404.txt` file inside
the folder "statuses".

- Status Code Tracking: Catalogues URLs based on the HTTP status codes they return. Each status code has its
corresponding text file, providing a clear and organized overview of your site's link health.

- Preparation for Sitemap creation: Capable of retrieving all reachable URLs on a website. This feature can be
used to create a comprehensive sitemap, improving your website's navigability and search engine visibility. All 
correctly working URLs can be found in `200.txt` file inside the folder "statuses".

By utilizing Link Checker, you can ensure that your website is free of broken links, effectively manage your
site's status code responses, and enhance your website's structure with the creation of comprehensive sitemaps.
Maintaining these aspects will significantly contribute to providing an optimal user experience and boosting your
site's SEO performance.

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

## Usage
To use Link Checker, follow these steps:

1. Open your terminal or command prompt.

2. Navigate to the directory where you have the `linkchecker` executable.

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

4. The Link Checker will scan the webpage, extract links, and check the status of each link. Then
it will visit each link and extract all URLs present on these secondary level pages, and so on. 
Starting from the home page, it will visit all the available links on your website.

   The results will be displayed in the folder "statuses". For Windows users, this folder will
   be created automatically. For Linux users, it may be necessary to run the program as `root` or
   to manually create the "statuses" folder with `777` access rights assigned.
   ![Screenshot_1](https://github.com/luxuriant777/linkchecker/assets/20545475/b0a48a30-2fe4-42b9-bb30-fc4f71db58e4)

    Example output:
   - `200.txt`:
    ```shell
    https://example.com/
    https://example.com/page1
    https://example.com/page2
    ```
   - `404.txt`:
    ```shell
   https://example.com/ -> https://example.com/error-test
   https://example.com/page3 -> https://example.com/error-test-second
   https://example.com/page4 -> https://example.com/error-test

    ```
   Here, `https://example.com/` is the page where the broken link was identified, while `https://example.com/error-test`
   is the broken link itself. A single page may contain multiple broken links, all of which will be listed in the
   results.

   That's it! You can now use these results to effectively address the found broken links on your website or create
   comprehensive sitemaps using the results offered by the `200.txt` output.

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please create an issue 
or submit a pull request on the `dev` branch.

## License

This project is licensed under the
[GNU General Public License (GPL) version 3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

## Disclaimer
The Link Checker program is provided as-is, without any warranties or guarantees of any kind. The authors and
contributors of this program disclaim any liability for damages or losses that may arise from its use.

The tool is designed to assist in identifying broken links and generating URL lists based on status codes.
However, it falls under your obligation to use this tool judiciously and comply with the terms of service and
guidelines of the websites you are examining.

Please exercise caution and use the Link Checker at your own risk. We recommend testing it on a small
scale or non-production environment before using it extensively. 
