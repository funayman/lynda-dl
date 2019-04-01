# lynda-dl

This is a small utility for downloading lynda.com courses if you have an organizational account.
If you have a regular lynda account, I recommend using [youtube-dl](https://github.com/rg3/youtube-dl) or [lynder](https://github.com/EnesCakir/lynder) as they are more sophisticated applications.

The design is inspired by [lynder](https://github.com/EnesCakir/lynder)

## Examples
**REQUIRES CURL TO BE INSTALLED ON SYSTEM**

Make sure you export your cookies from your browser! There are multiple extensions for Chrome and Firefox that allow you to do this.

### Default
```
Usage:
  lynda-dl [command]

Available Commands:
  download    Download a Lynda course
  help        Help about any command
  list        view the contents of a course

Flags:
      --config string   config file (default is $HOME/.lynda-dl.yaml)
  -i, --course-id int   Lynda course id
  -h, --help            help for lynda-dl

Use "lynda-dl [command] --help" for more information about a command.
```

### Download Course
```
$ lynda-dl download -c ~/Downloads/cookies.txt https://www.lynda.com/Go-tutorials/Go-Essential-Training/748576-2.html
```
### Download a Learning Path
```
$ lynda-dl download -c ~/Downloads/cookies.txt --learning-path https://www.lynda.com/learning-paths/Audio-Music/become-a-music-producer
```
### List the Contents of a Course
```
$ lynda-dl list https://www.lynda.com/Developer-tutorials/Agile-Software-Development-Code-Quality/5005070-2.html
Agile Software Development: Code Quality
|--> Introduction
|  |--> Why code quality is important
|  |--> What you should know
|--> 1. Code Quality
|  |--> What is code quality?
|  |--> How do we end up with poor quality code?
|  |--> Review of code quality metrics and tools
|--> 2. Complexity
|  |--> What is complexity?
|  |--> Exploring complexity with Code Climate
|  |--> Enabling cyclomatic complexity
|  |--> Adjust thresholds
|--> 3. Hotspots and Churn
|  |--> What are hotspots and churn?
|  |--> Hotspots and churn with CodeScene
|  |--> Customizing analysis with CodeScene
|--> 4. Code Coverage
|  |--> What is code coverage?
|  |--> Windows setup
|  |--> Collecting code coverage with dotCover
|  |--> Visualizing code coverage with NDepend
|  |--> macOS X setup
|  |--> Collecting code coverage with SimpleCov
|  |--> Visualizing code coverage with Code Climate
|--> 5. Duplication
|  |--> What is duplication?
|  |--> Setting up copy/paste detectors (CPD)
|  |--> Finding duplication with CPD
|  |--> Visualizing duplication with SonarQube
|--> 6. Securing Your Dependencies
|  |--> Why secure dependencies?
|  |--> Setting up Snyk
|  |--> Securing your dependencies with Snyk CLI
|  |--> Securing your dependencies with Snyk web
|--> 7. Consistent Coding Style
|  |--> Why consistent coding style?
|  |--> Setting up Flask
|  |--> Detect style violations with Flake8
|  |--> Automated code style review with Hound CI
|--> Conclusion
|  |--> Next steps
```
