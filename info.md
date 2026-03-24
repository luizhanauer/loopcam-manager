#

##
/home/unk/Vídeos/loop-cam/seu_video.mp4

##
1. Remove o módulo atual
sudo rmmod v4l2loopback

2. Recarrega com exclusive_caps=1
sudo modprobe v4l2loopback video_nr=10 card_label="LoopCam-Browser" exclusive_caps=1

##

sudo fuser -k /dev/video10
sudo rmmod v4l2loopback
sudo modprobe v4l2loopback video_nr=10 card_label="LoopCam-Browser" exclusive_caps=1