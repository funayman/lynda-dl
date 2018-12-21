# lynda-dl

**REQUIRES CURL TO BE INSTALLED ON SYSTEM**

This is a small utility for downloading lynda.com courses if you have an organizational account.
If you have a regular lynda account, I recommend using [youtube-dl](https://github.com/rg3/youtube-dl) or [lynder](https://github.com/EnesCakir/lynder) as they are more sophisticated applications.

The design is inspired by [lynder](https://github.com/EnesCakir/lynder)

## Usage
```bash
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

```bash
lynda-dl download -c ~/Downloads/cookies.txt --course-id 573393
```

Make sure you export your cookies from your browser. There are multiple extensions for Chrome and Firefox that allow you to do this.

