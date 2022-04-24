# eksdemo
instalation guide
download the source to your local machine -- https://github.com/aaroniscode/eksdemo
eksdemo requires brew as the package manager. You can download this package manager here - https://brew.sh/
Downlaod and install homebrew - /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
Set path as follows
test -d ~/.linuxbrew && eval $(~/.linuxbrew/bin/brew shellenv)
test -d /home/linuxbrew/.linuxbrew && eval $(/home/linuxbrew/.linuxbrew/bin/brew shellenv)
test -r ~/.bash_profile && echo "eval \$($(brew --prefix)/bin/brew shellenv)" >>~/.bash_profile
echo "eval \$($(brew --prefix)/bin/brew shellenv)" >>~/.profile
check version as follows -- brew --version
![image](https://user-images.githubusercontent.com/12427542/164969520-b52d1147-6269-43ed-8290-b08bcc937047.png)
check if eksdemo is working -- eksdemo --help
![image](https://user-images.githubusercontent.com/12427542/164969600-53daf43f-c80d-4631-ae26-a4754b68845b.png)
There you go!!!
