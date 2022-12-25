
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

# Install wine
ARG WINE_BRANCH="stable"
RUN wget -nv -O- https://dl.winehq.org/wine-builds/winehq.key | APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=1 apt-key add - \
    && apt-add-repository "deb https://dl.winehq.org/wine-builds/ubuntu/ $(grep VERSION_CODENAME= /etc/os-release | cut -d= -f2) main" \
    && dpkg --add-architecture i386 \
    && apt-get update \
    && DEBIAN_FRONTEND="noninteractive" apt-get install -y --install-recommends winehq-${WINE_BRANCH} \
    && rm -rf /var/lib/apt/lists/*

# Install winetricks
RUN wget -nv -O /usr/bin/winetricks https://raw.githubusercontent.com/Winetricks/winetricks/master/src/winetricks \
    && chmod +x /usr/bin/winetricks

# Download gecko and mono installers
COPY download_gecko_and_mono.sh /root/download_gecko_and_mono.sh
RUN chmod +x /root/download_gecko_and_mono.sh \
    && /root/download_gecko_and_mono.sh "$(dpkg -s wine-${WINE_BRANCH} | grep "^Version:\s" | awk '{print $2}' | sed -E 's/~.*$//')"

# Install graphic libraries
RUN winetricks d3dx9_43
# RUN winetricks --force -q dotnet48

# Silence all fixme warnings from wine
ENV WINEDEBUG fixme-all

WORKDIR /appvm

COPY default.pa /etc/pulse/
COPY supervisord.conf /etc/supervisor/conf.d/
COPY syncinput.cpp ./syncinput.cpp

# Compile syncinput.cpp
RUN x86_64-w64-mingw32-g++ ./syncinput.cpp -o ./syncinput.exe -lws2_32 -lpthread -static

# Add reg files to wineprefix
ARG APP_NAME
COPY apps/${APP_NAME}/setup/ /root/.wine/

RUN add-apt-repository ppa:lutris-team/lutris
RUN apt-get update

COPY lutris/lutris_0.5.12_all.deb lutris_0.5.12_all.deb
RUN dpkg -i lutris_0.5.12_all.deb; exit 0
RUN DEBIAN_FRONTEND=noninteractive apt --fix-broken install -y
RUN DEBIAN_FRONTEND=noninteractive apt install lutris -y

ENTRYPOINT ["supervisord"]