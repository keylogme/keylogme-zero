#!/bin/bash
# variables
keylogme_version=$1
file_config_abs_path=$2
# Check sudo permissions
echo $EUID
if [ $EUID != 0 ]; then
    echo "🟡 Running installer with non-sudo permissions."
    echo "   Please run the script with sudo privileges to create keylogger service"
    echo ""
    exit 1
fi

# Check whether the given command exists.
has_cmd() {
    command -v "$1" > /dev/null 2>&1
}

has_cmd wget || {
    echo "🟡 wget is not installed. Please install wget to download keylogme-zero"
    exit 1
}

has_cmd tar || {
    echo "🟡 tar is not installed. Please install tar to extract keylogme-zero"
    exit 1
}

has_cmd systemctl || {
    echo "🟡 systemctl is not installed. Please install systemctl to create keylogger service"
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

echo "Downloading keylogme-zero ${keylogme_version}..."

# download
file_compressed="keylogme-zero_Linux_x86_64.tar.gz"
wget -q https://github.com/keylogme/keylogme-zero/releases/download/${keylogme_version}/${file_compressed} -O ${file_compressed}

# unzip
echo "Uncompressing keylogme-zero ${keylogme_version}..."
mkdir keylogme
tar -xvzf ${file_compressed} -C keylogme

# check if service keylogme-zero exists and stop it
systemctl is-active --quiet keylogger-zero && {
    echo "🟡 keylogger-zero service is running. Stopping the service..."
    sudo systemctl stop keylogger-zero
}

# try to copy and check if failed
cp keylogme/keylogme-zero /bin || {
    echo "🟡 Failed to copy keylogme-zero to /bin"
    exit 1
}

# Copy service file, incase if there are any changes
service_file_path="/etc/systemd/system/keylogger-zero.service"
sudo cp keylogger-zero.service ${service_file_path}
# Set environment variables in service file
echo $"Environment=CONFIG_FILE=${file_config_abs_path}" >> ${service_file_path}
# reload configurations incase if service file has changed
sudo systemctl daemon-reload
# restart the service
sudo systemctl restart keylogger-zero
# start of VM restart
sudo systemctl enable keylogger-zero

# check service keylogme-zero is running
systemctl is-active --quiet keylogger-zero && {
    echo "🟢 keylogger-zero service is running."
}

echo "🟢 keylogme-zero ${keylogme_version} installed successfully"
