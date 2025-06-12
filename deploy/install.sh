#!/bin/bash

set -o errexit

echo "##############################################################################"
cat << 'EOF'



 _              _                                                   
| | _____ _   _| | ___   __ _ _ __ ___   ___      _______ _ __ ___  
| |/ / _ \ | | | |/ _ \ / _` | '_ ` _ \ / _ \____|_  / _ \ '__/ _ \ 
|   <  __/ |_| | | (_) | (_| | | | | | |  __/_____/ /  __/ | | (_) |
|_|\_\___|\__, |_|\___/ \__, |_| |_| |_|\___|    /___\___|_|  \___/ 
          |___/         |___/                                       



EOF


echo "##############################################################################"
echo ""
echo "PRE STEPS"
echo ""
echo "##############################################################################"

echo "Running with shell: $SHELL"
echo "Current interpreter: $(ps -p $$ -o comm=)" # $$ is current PID, comm= is command name
#
# variables
keylogme_version=$1
file_config_abs_path=$2
# Check sudo permissions
if [ $EUID != 0 ]; then
    echo "🟡 Running installer with non-sudo permissions."
    echo "   Please run the script with sudo privileges to create keylogme service"
    echo ""
    exit 1
fi

has_cmd() {
    command -v "$1" > /dev/null 2>&1
}

# Check whether 'wget' command exists.
has_wget() {
    has_cmd wget
}

# Check whether 'curl' command exists.
has_curl() {
    has_cmd curl
}

has_tar(){
    has_cmd tar
}

has_systemctl(){
    has_cmd systemctl
}

has_launchctl(){
    has_cmd launchctl
}

is_mac() {
    [[ `uname -s` == 'Darwin' || `uname -s` == 'MacOS' || `uname -s` == 'macOS' ]]
}

is_linux(){
    [[ `uname -s` == 'Linux' ]]
}

is_arm64(){
    [[ `uname -m` == 'arm64' || `uname -m` == 'aarch64' ]]
}

check_os(){
    if is_mac; then
        desired_os=1
        os="Darwin"
        return
    elif is_linux; then
        desired_os=1
        os="Linux"
        return
    else
        echo "🟡 Unsupported OS. Please run this script on Linux or macOS."
        exit 1
    fi
}

check_arch(){
    if is_arm64; then
        arch="arm64"
    else
        arch="x86_64"
    fi
}

# Function to pause script execution and wait for user to press Enter
function press_enter_to_continue() {
    echo ""
    read -p "Press Enter to continue..."
    echo ""
}

# Function to open System Settings to the Privacy & Security pane
# Works for macOS Ventura (13.0) and newer
function open_privacy_settings() {
    echo "Opening 'System Settings' to 'Privacy & Security'..."
    if sw_vers -productVersion | grep -E '^(1[3-9]|[2-9][0-9])' > /dev/null; then
        # macOS Ventura (13.0) and later
        open "x-apple.systempreferences:com.apple.preference.security"
    else
        # macOS Monterey (12.x) and older
        open "/System/Library/PreferencePanes/Security.prefPane"
    fi
}

# Function to open a Finder window to a specified folder
# Argument 1: The path to the desired folder
function open_folder_in_finder() {
    local folder_path="$1"

    # Check if a folder path was provided
    if [ -z "$folder_path" ]; then
        echo "Error: No folder path provided."
        echo "Usage: open_folder_in_finder /path/to/your/folder"
        return 1 # Indicate an error
    fi

    # Check if the provided path is a valid directory
    if [ ! -d "$folder_path" ]; then
        echo "Error: Directory '$folder_path' does not exist or is not a directory."
        return 1 # Indicate an error
    fi

    # Use the 'open' command with the specified folder path
    # 'open .' will open the current directory, 'open /path/to/folder' will open that folder.
    open "$folder_path"

    echo "Opened Finder window to: '$folder_path'"
    return 0 # Indicate success
}

# Check whether the given command exists.
has_cmd tar || {
    echo "🟡 tar is not installed. Please install tar to extract keylogme-zero"
    exit 1
}
has_cmd envsubst || {
    echo "🟡 envsubst is not installed. Please install envsubst to set environment variables in service file"
    exit 1
}
#######################
# Check OS
#######################
desired_os=0
os=""
arch=""
echo -e "🌏 Detecting your OS ..."
check_os
check_arch
echo "  OS ${os} arch ${arch}"

github_repo="keylogme/keylogme-zero"
keylogme_folder="${HOME}/.keylogme"
mkdir -p "$keylogme_folder"

# Check inputs
# Required
# (none)

# Optional
if [[ "$keylogme_version" == "" || "$keylogme_version" == "latest" ]]; then
    latest_release_info=$(curl -s "https://api.github.com/repos/$github_repo/releases/latest")

    # Extract the tag name
    keylogme_version=$(echo "$latest_release_info" | jq -r '.tag_name')
    echo "Latest keylogme version ${keylogme_version}"
fi


if [ "$file_config_abs_path" == "" ]; then
    dir_name="$(pwd)"
    if [ "$os" == "Linux" ] ;then
        default_name="default_config_linux.json.template"
    elif [ "$os" == "Darwin" ]; then
        default_name="default_config_darwin.json.template"
    fi

    file_config_abs_path="${dir_name}/${default_name}"
    echo "Using template config file path"

    output_file="${keylogme_folder}/output-keylogme-zero.json"
    export KEYLOGME_OUTPUT_FILE="${output_file}"
    envsubst < "${file_config_abs_path}" > "${keylogme_folder}/config-keylogme-zero.json"
    file_config_abs_path="${keylogme_folder}/config-keylogme-zero.json"
    echo "##############################################################################"
    echo ""
    echo "☀️Here are the interesting files: the config and ouput"
    echo ""
    echo "Output file will be saved to ${output_file}"
    echo "Config file will be saved to ${file_config_abs_path}"
    echo ""
    echo ""
    echo "##############################################################################"
fi

echo ""
echo "START OF INSTALLATION"
echo ""



# download
echo "##############################################################################"
echo "Step 1: ⬇️Downloading keylogme-zero ${keylogme_version}..."
file_compressed="keylogme-zero_${os}_${arch}.tar.gz"
url="https://github.com/${github_repo}/releases/download/${keylogme_version}/${file_compressed}"
echo "  File to download: ${url}"
if has_curl; then
    echo "🟢 Using curl to download keylogme-zero..."
    curl -s -L ${url} --output ${file_compressed} 
elif has_wget; then
    echo "🟢 Using wget to download keylogme-zero..."
    wget -q ${url} -O ${file_compressed}
else
    echo "🟡 No download tool found. Please install curl or wget or fetch to download keylogme-zero."
    exit 1
fi

# unzip
echo "##############################################################################"
echo "Step 2: 🗜️Uncompressing keylogme-zero ${keylogme_version}..."
mkdir -p keylogme
if has_tar; then
    tar -xzf ${file_compressed} -C keylogme
else
    echo "🟡 tar command not found. Please install tar to extract keylogme-zero."
    exit 1
fi

# check if service keylogme-zero exists and stop it
echo "##############################################################################"
echo "Step 3: 🍏Checking if keylogme-zero service exists..."
if [ "$os" == "Linux" ] ;then
    service_file_path="/etc/systemd/system/keylogme-zero.service"
elif [ "$os" == "Darwin" ]; then
    service_file_path="/Library/LaunchDaemons/com.keylogme.keylogme-zero.plist"
fi

if has_systemctl; then
    systemctl is-active --quiet keylogme-zero && {
        echo "🟡 keylogme-zero service is running. Stopping the service..."
        sudo systemctl stop keylogme-zero
        sudo systemctl disable keylogme-zero
        echo "  Service was stopped"
    }
elif has_launchctl; then
    launchctl list | grep -q keylogme-zero && {
        echo "🟡 keylogme-zero service is running. Stopping the service..."
        sudo launchctl stop com.keylogme.keylogme-zero
        sudo launchctl unload ${service_file_path}
        echo "  Service was stopped"
    }
else
    echo "🟡 Neither systemctl(Linux) nor launchctl(MacOS) command found. Please run this script on a system with systemd or launchd."
    exit 1
fi


echo "##############################################################################"
echo "Step 4: 🎯Copy binary..."
binary_folder=""
if [ "$os" == "Linux" ] ;then
    # copy binary
    binary_folder="/bin/"
elif [ "$os" == "Darwin" ]; then
    binary_folder="/usr/local"
    if [ "$arch" == "arm64" ]; then
        binary_folder="/opt"
    fi
fi
keylogme_version_fmt="${keylogme_version//./_}"
binary_name="keylogme-zero_${keylogme_version_fmt}"
binary_full_path="${binary_folder}/${binary_name}"

sudo cp ./keylogme/keylogme-zero "${binary_full_path}" || {
    echo "🟡 Failed to copy to ${binary_full_path}"
    exit 1
}


echo "##############################################################################"
echo "Step 5: 🚓Grant permissions..."
if [ "$os" == "Linux" ] ;then
    echo "You are running with sudo. Go ahead."
elif [ "$os" == "Darwin" ]; then
    echo ""
    echo "Manually grant 'Input Monitoring' permission in System Settings to ${binary_folder}/keylogme-zero"
    echo ""
    echo "Why does it need 'Input Monitoring'?"
    echo "keylogme-zero tracks usage per device. Input monitoring permission is the only"
    echo "way to know if a key was pressed from your built-in keyboard or an external keyboard."
    echo ""
    echo "Instructions:"
    echo "  a. Go to 'System Settings' (or 'System Preferences' on older macOS)."
    echo "  b. Click on 'Privacy & Security' in the sidebar."
    echo "  c. Scroll down and click on 'Input Monitoring'."
    echo "  d. In another Finder window. From the toolbar > Go > Go to folder and type ${binary_folder}"
    echo "  e. Drag the binary ${binary_name} to the 'Input Monitoring' allowed applications (make sure it is enabled)"
    echo "  WARNING: if you are reinstalling, you have to drag again!"
    echo ""
    echo ""
    # Ask user if they want to open settings
    read -p "Would you like me to open 'System Settings->Privacy & Security' and the binary folder for you now? (y/n): " response_open_settings
    if [[ "$response_open_settings" =~ ^[Yy]$ ]]; then
        open_privacy_settings
        open_folder_in_finder "${binary_folder}"
    fi
    press_enter_to_continue
fi



echo "##############################################################################"
echo "Step 6: 🖥️Setting up service keylogme-zero..."
export KEYLOGME_ZERO_CONFIG_FILE_PLACEHOLDER=${file_config_abs_path}
export KEYLOGME_BINARY_PLACEHOLDER="${binary_full_path}"
if [ "$os" == "Linux" ] ;then
    envsubst < keylogme-zero.service.template > ${service_file_path}
    # Set environment variables in service file
    # echo $"Environment=CONFIG_FILE=${file_config_abs_path}" >> ${service_file_path}
    # reload configurations incase if service file has changed
    sudo systemctl daemon-reload
    sudo systemctl enable keylogme-zero
    sudo systemctl restart keylogme-zero
    # check service keylogme-zero is running
    systemctl is-active --quiet keylogme-zero && {
        echo "🟢 keylogme-zero service is running."
    }
elif [ "$os" == "Darwin" ]; then
    envsubst < keylogme-zero.plist.template > ${service_file_path}

    sudo chown root:wheel ${service_file_path}
    sudo chmod 644 ${service_file_path}

    sudo launchctl load -w ${service_file_path}
    sudo launchctl start com.keylogme.keylogme-zero
fi


# check service is running
echo "##############################################################################"
echo "Step 7: 🪄Checking service is running..."
sleep 5

if has_systemctl; then
    systemctl is-active --quiet keylogme-zero || {
        echo "❌ Service is not running"
        exit 1
    }
    echo "🆗 Service is running"
elif has_launchctl; then
    launchctl list | grep -q keylogme-zero || {
        echo "❌ Service is not running"
        exit 1
    }
    echo "🆗 Service is running"
else
    echo "🟡 Neither systemctl(Linux) nor launchctl(MacOS) command found. Please run this script on a system with systemd or launchd."
    exit 1
fi

echo ""
echo ""
echo "🟢 keylogme-zero ${keylogme_version} installed successfully"
echo ""
echo "What is doing? check logs by running command below:"
echo ""
if [ "$os" == "Linux" ] ;then
    echo "systemctl status keylogme-zero"
elif [ "$os" == "Darwin" ]; then
    echo "cat /var/log/keylogme-zero.log"
fi
echo ""
echo "Next steps?"
echo ""
echo "1. Modify the config file. See ref : https://github.com/keylogme/keylogme-zero?tab=readme-ov-file#config"
echo "2. If you modify config file, restart service to update and use your config"
if [ "$os" == "Linux" ] ;then
    echo "   sudo systemctl restart keylogme-zero"
elif [ "$os" == "Darwin" ]; then
    echo "   sudo launchctl stop com.keylogme.keylogme-zero"
    echo "   sudo launchctl start com.keylogme.keylogme-zero"
fi
echo ""
echo ""
echo "You can visualize your keyboard usage (heatmap) in https://keylogme/anonymous"
echo ""
echo "🌅 Have a coconut oil smooth typing 🌴"
echo ""
echo ""
