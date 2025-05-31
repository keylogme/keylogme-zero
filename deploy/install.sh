#!/bin/bash

set -o errexit

echo "Running with shell: $SHELL"
echo "Current interpreter: $(ps -p $$ -o comm=)" # $$ is current PID, comm= is command name
#
# variables
keylogme_version=$1
file_config_abs_path=$2
# Check sudo permissions
echo $EUID
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

# Check whether the given command exists.
has_cmd tar || {
    echo "🟡 tar is not installed. Please install tar to extract keylogme-zero"
    exit 1
}
has_cmd envsubst || {
    echo "🟡 envsubst is not installed. Please install envsubst to set environment variables in service file"
    exit 1
}

# Check inputs
# Required
# (none)

# Optional
if [ "$keylogme_version" == "" ]; then
    echo "keylogme version default to latest"
    keylogme_version="latest"
fi
if [ "$file_config_abs_path" == "" ]; then
    dir_name="$(pwd)"
    default_name="default_config.json"
    file_config_abs_path="${dir_name}/${default_name}"
    echo "Absolute config file path will be set to ${file_config_abs_path}"
fi

#############################################################################
# START OF INSTALLATION
#############################################################################


# check OS
desired_os=0
os=""
arch=""
echo -e "🌏 Detecting your OS ...\n"
check_os
check_arch

# download
echo "⬇️Downloading keylogme-zero ${keylogme_version}..."
file_compressed="keylogme-zero_${os}_${arch}.tar.gz"
url="https://github.com/keylogme/keylogme-zero/releases/download/${keylogme_version}/${file_compressed}"
echo "  File to download: ${url}"
if has_curl; then
    echo "🟢 Using curl to download keylogme-zero..."
    curl -v -L ${url} --output ${file_compressed} 
elif has_wget; then
    echo "🟢 Using wget to download keylogme-zero..."
    wget -q ${url} -O ${file_compressed}
else
    echo "🟡 No download tool found. Please install curl or wget or fetch to download keylogme-zero."
    exit 1
fi

# unzip
echo "🗜️Uncompressing keylogme-zero ${keylogme_version}..."
mkdir -p keylogme
if has_tar; then
    tar -xvzf ${file_compressed} -C keylogme
else
    echo "🟡 tar command not found. Please install tar to extract keylogme-zero."
    exit 1
fi

# check if service keylogme-zero exists and stop it
echo "🍏Checking if keylogme-zero service exists..."
if os="Linux";then
    service_file_path="/etc/systemd/system/keylogme-zero.service"
elif os="Darwin"; then
    service_file_path="/Library/LaunchDaemons/keylogme-zero.plist"
fi

if has_systemctl; then
    systemctl is-active --quiet keylogme-zero && {
        echo "🟡 keylogme-zero service is running. Stopping the service..."
        sudo systemctl stop keylogme-zero
        sudo systemctl disable keylogme-zero
    }
elif has_launchctl; then
    launchctl list | grep -q keylogme-zero && {
        echo "🟡 keylogme-zero service is running. Stopping the service..."
        sudo launchctl stop com.keylogme.keylogme-zero
        sudo launchctl unload ${service_file_path}
else
    echo "🟡 Neither systemctl(Linux) nor launchctl(MacOS) command found. Please run this script on a system with systemd or launchd."
    exit 1

fi


# try to copy and check if failed
sudo cp keylogme/keylogme-zero /bin || {
    echo "🟡 Failed to copy keylogme-zero to /bin"
    exit 1
}

export KEYLOGME_ZERO_CONFIG_FILE_PLACEHOLDER=${file_config_abs_path}
envsubst < keylogme/keylogme-zero.service.template > keylogme-zero.service 
if os="Linux"; then
    sudo cp keylogme-zero.service ${service_file_path}
    # Set environment variables in service file
    # echo $"Environment=CONFIG_FILE=${file_config_abs_path}" >> ${service_file_path}
    # reload configurations incase if service file has changed
    sudo systemctl daemon-reload
    # restart the service
    sudo systemctl restart keylogme-zero
    # start of VM restart
    sudo systemctl enable keylogme-zero
    # check service keylogme-zero is running
    systemctl is-active --quiet keylogme-zero && {
        echo "🟢 keylogme-zero service is running."
    }
elif os="Darwin"; then
    sudo cp keylogme-zero.plist ${service_file_path}
fi


echo "🟢 keylogme-zero ${keylogme_version} installed successfully"
