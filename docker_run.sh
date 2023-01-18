app=/home/owais/CloudGaming/Selkies/CopeGaming/appvm/apps/gtavc_pc
videoport=37272
audioport=50405
wsport=46715

envfile=/home/owais/CloudGaming/Selkies/CopeGaming/provider/appconf/gtavc_pc.env

docker run -it -v $app:/appvm/app -e videoport=$videoport -e audioport=$audioport -e wsport=$wsport --env-file=$envfile --rm --gpus='all,"capabilities=compute,utility,graphics,display,video"' provider_appvm:latest