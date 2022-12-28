echo "Configure Lutris.... \n"

rm -r ../appvm/lutris/conf/*
rm -r ../appvm/lutris/share/*

cp -r /home/owais/.config/lutris/games/ ../appvm/lutris/config/
cp -r /home/owais/.local/share/lutris/ ../appvm/lutris/share/
