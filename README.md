<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/keylogme/keylogme-zero">
    <img src="images/icon-keylogme-zero-80x80.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Keylogme Zero</h3>

  <p align="center">
    This is a logger for <a href="https://keylogme.com">keylogme.com</a> . This logger saves 
    the stats locally. You can use those results to visualize in (pending...).
    <br />
    <br />
    <a href="https://keylogme.com/esoteloferry">View Demo</a>
    ·
    <a href="https://github.com/keylogme/keylogme-zero/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/keylogme/keylogme-zero/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#security">Security</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#linux">Linux</a></li>
      </ul>
    </li>
    <li><a href="#config">Config</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

The end goal is to avoid or diminish the pain in your hands due to typing. Commonly known as RSI (Repetitive Strain Injury), 
tennis elbow, carpal tunnel syndrome, etc.

There are many great ergonomic keyboards , beautiful hardware out there. However, what is the layout of keys 
that best suits you? that does not cause you pain? that makes you more productive?.

We all started with QWERTY, then heard of DVORAK, COLEMAK, Workman, Norman, Asset, Capewell-Dvorak, BEAKL, MTGAP, QGMLWB... ?
There are many layouts but switching to one is not an easy task, you need a lot of practice and patience.


Here's how:
* Monitor : See the finger usage on your layout based on your real usage
* Analyze : Compare your layout with others, find patterns to avoid or improve, remap shortcuts
* Adapt : fine tune your layout based on the stats

Of course, ergonomics is not just a nice keyboard and layout. It is also about posture, breaks, exercises.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Security

A keylogger is a tool that records the keystrokes on a computer. It can be used for good or bad purposes.
Of course, our intention is to use it for good purposes. How can you trust that?, well the code is completely open source, 
no dependencies and it stores your data locally in your computer, there is no connection to the internet.

The online viewer does not need an account to use it. You can use it anonymously to visualize your stats. 

<!-- GETTING STARTED -->
## Getting Started

### Linux

1. Clone the repo
   ```sh
   git clone https://github.com/keylogme/keylogme-zero.git
   ```
2. Go to deploy and install with sudo permissions
   ```sh
   cd deploy && sudo ./install.sh
   ```
<details>
  <summary>With parameters</summary>
   If you want to install a specific version:
   ```sh
   cd deploy && sudo ./install.sh v1.1.0
   ```
   If you want to install and use your own config (don't forget the version, in this case latest):
   ```sh
   cd deploy && sudo ./install.sh latest /path/to/your/config.json
   ```
</details>

3. After some keypresses and 10 seconds, check the stats in `/output_keylogme_zero.json`


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- Config EXAMPLES -->
## Config

The default config in `deploy` folder is:

```json
{
    "keylog": {
        "devices": [
            {
                "device_id": "1",
                "name": "✋ crkbd ⌨️",
                "usb_name": "foostan Corne"
            }
        ],
        "shortcuts": [
            {
                "id": 1,
                "codes": [
                    29,
                    46
                ],
                "type": "hold"
            }
        ]
    },
    "storage": {
        "file_output": "output_keylogme_zero.json",
        "periodic_save_in_sec": 10
    }
}
```

The config has two main sections:

- keylog : config for keylogger
    - devices : list of devices to monitor
        - device_id : unique id for the device
        - name : name of the device, you named it as you want.
        - usb_name : usb name of the device. Go to <a href="#usb-name">USB name</a> section to know how to get it.
    - shortcuts : list of shortcuts to monitor
        - id : unique id for the shortcut
        - codes : list of keycodes (decimal format) for the shortcut. Go to <a href="#keycodes-hardware">Keycodes hardware</a> section to know how to get it.
        - type : type of shortcut (hold, press)
- storage : config for storage
    - file_output : abs filepath to store the stats
    - periodic_save_in_sec : periodic time to save the stats. In seconds.


<p align="right">(<a href="#readme-top">back to top</a>)</p>

### USB name

A usb device connected to computer has a unique name. 
To get the name of the device, you can use the following command:

```sh
apt install input-utils
sudo lsinput
```

If your keyboard name appeared multiple times, try with all of them.

For example, the output of the command is below, and the name that worked is `foostan Corne`.

```sh
/dev/input/event12
   bustype : BUS_USB
   vendor  : 0x4653
   product : 0x1
   version : 273
   name    : "foostan Corne"
   phys    : "usb-0000:00:14.0-4.3/input0"
   uniq    : ""
   bits ev : (null) (null) (null) (null) (null)

/dev/input/event13
   bustype : BUS_USB
   vendor  : 0x4653
   product : 0x1
   version : 273
   name    : "foostan Corne System Control"
   phys    : "usb-0000:00:14.0-4.3/input2"
   uniq    : ""
   bits ev : (null) (null) (null) (null)

/dev/input/event14
   bustype : BUS_USB
   vendor  : 0x4653
   product : 0x1
   version : 273
   name    : "foostan Corne Consumer Control"
   phys    : "usb-0000:00:14.0-4.3/input2"
   uniq    : ""
   bits ev : (null) (null) (null) (null) (null)

/dev/input/event15
   bustype : BUS_USB
   vendor  : 0x4653
   product : 0x1
   version : 273
   name    : "foostan Corne Keyboard"
   phys    : "usb-0000:00:14.0-4.3/input2"
   uniq    : ""
   bits ev : (null) (null) (null) (null) (null)
```


### Keycodes hardware

A key(hardware) has a keycode, f.e. in a normal QWERTY keyboard, the keycode of Q is 
10(HEX) and 16(Decimal), letter C is 2E(HEX) and 46(Decimal). 

The keyboard (hardware) sends the keycode to the computer. The computer uses the
keyboard layout to convert the keycode to a character. The keyboard layout is defined 
in your operating system. For example, the layout US QWERTY will convert 16(Decimal) to Q 
and 46(Decimal) to C. But if I have defined another layout, for example 
[WORKMAN](https://workmanlayout.org/), then Q 16(Decimal)
will be Q and 46(Decimal) will be M. You get the idea 🙃

How to know the keycode?
Go to [kbdlayout.info](https://kbdlayout.info/kbdus)
and click scancodes to see the keycodes. 
The scancode is a hex number, you have to convert it to decimal.

<!-- ROADMAP -->
## Roadmap


TODO
- [x] Add Changelog
- [x] Add back to top links
- [ ] Add Additional Templates w/ Examples
- [ ] Add "components" document to easily copy & paste sections of the readme
- [ ] Multi-language Support
    - [ ] Chinese
    - [ ] Spanish

See the [open issues](https://github.com/othneildrew/Best-README-Template/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request




<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Efrain Sotelo - [@esoteloferry](https://twitter.com/esoteloferry)

Project Link: [https://keylogme.com](https://keylogme.com)

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

Helpful resources and would like to give credit to. 

* [Linux keylogger](https://github.com/MarinX/keylogger) helpful starting point

<p align="right">(<a href="#readme-top">back to top</a>)</p>


