
ARG UBUNTU_RELEASE=20.04
ARG CUDA_VERSION=11.2.2
FROM nvcr.io/nvidia/cudagl:${CUDA_VERSION}-runtime-ubuntu${UBUNTU_RELEASE}

ARG UBUNTU_RELEASE
ARG CUDA_VERSION
# Make all NVIDIA GPUs visible by default
ARG NVIDIA_VISIBLE_DEVICES=all
# Use noninteractive mode to skip confirmation when installing packages
ARG DEBIAN_FRONTEND=noninteractive
# All NVIDIA driver capabilities should preferably be used, check `NVIDIA_DRIVER_CAPABILITIES` inside the container if things do not work
ENV NVIDIA_DRIVER_CAPABILITIES all

# Install dependencies
RUN apt update \
    && apt-get update -y \
    && apt-get install --no-install-recommends --assume-yes \
    wget \
    software-properties-common \
    gpg-agent \
    supervisor \
    xvfb \
    mingw-w64 \
    ffmpeg \
    cabextract \
    aptitude \
    vim \
    pulseaudio \
    mesa-utils \
    mesa-utils-extra \
    mesa-va-drivers \
    mesa-vulkan-drivers \
    libvulkan-dev \
    libvulkan-dev:i386 \
    && apt-get clean \
    && apt-get autoremove


# Install and configure Vulkan manually
RUN if [ "${UBUNTU_RELEASE}" = "18.04" ]; then apt-get update && apt-get install --no-install-recommends -y vulkan-utils; else apt-get update && apt-get install --no-install-recommends -y vulkan-tools; fi && \
    rm -rf /var/lib/apt/lists/* && \
    VULKAN_API_VERSION=$(dpkg -s libvulkan1 | grep -oP 'Version: [0-9|\.]+' | grep -oP '[0-9]+(\.[0-9]+)(\.[0-9]+)') && \
    mkdir -p /etc/vulkan/icd.d/ && \
    echo "{\n\
    \"file_format_version\" : \"1.0.0\",\n\
    \"ICD\": {\n\
        \"library_path\": \"libGLX_nvidia.so.0\",\n\
        \"api_version\" : \"${VULKAN_API_VERSION}\"\n\
    }\n\
}" > /etc/vulkan/icd.d/nvidia_icd.json

# Silence all fixme warnings from wine
ENV WINEDEBUG fixme-all

# Install Lutris dependencies
# Add 32 bit architecture and install wine32
RUN dpkg --add-architecture i386 && apt-get update -y

COPY Lutris-Deps /Lutris-Deps
WORKDIR /Lutris-Deps
RUN dpkg -i ./lutris-build-deps_0.5.12-0ubuntu1ppa1~ubuntu22.04-0ubuntu1ppa1~ubuntu22.04-0ubuntu1ppa1~ubuntu22.04_all.deb; \
    apt-get --fix-broken install -y; apt-get update --fix-missing; 

# Clone Lutris-Source
RUN apt install git -y
RUN git clone --branch NoRootWarning https://github.com/joyflix99/lutris.git;

RUN apt install wine32 -y

RUN apt-get install -y dbus-x11 libcanberra-gtk-module libcanberra-gtk3-module pciutils

WORKDIR /appvm

COPY default.pa /etc/pulse/
COPY supervisord.conf /etc/supervisor/conf.d/
COPY syncinput.cpp ./syncinput.cpp

# Compile syncinput.cpp
RUN x86_64-w64-mingw32-g++ ./syncinput.cpp -o ./syncinput.exe -lws2_32 -lpthread -static

# Add reg files to wineprefix
ARG APP_NAME
COPY apps/${APP_NAME}/setup/ /root/.wine/

# Configure Lutris
RUN apt install sqlite3
COPY lutris/share/ /root/.local/share/lutris/
COPY lutris/config/ /root/.config/lutris/

ENTRYPOINT ["supervisord"]