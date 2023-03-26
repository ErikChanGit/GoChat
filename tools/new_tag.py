import subprocess
# git config --global user.email "1043490933@qq.com"
#             git config --global user.name "ErikChanGit"
#             git tag -a v0.1.8 -m "Release_v0.1.8"
#             git push origin v0.1.8


# cd ${{ github.workspace }}/tools
#             python new_tag.py

def clone_from_git(repository, bsp_repository_dest):
    ''' 将 github 的源码克隆到本地 '''
    if os.path.exists(bsp_repository_dest + "/.git") is False:
        subprocess.check_call("git clone " + repository + " " + bsp_repository_dest)


if __name__ == '__main__':
    clone_from_git("git@github.com:ErikChanGit/GoChat.git", "/sdk/GoChat")
    subprocess.check_call('git tag -a v0.1.8 -m "Release_v0.1.8"', cwd="/sdk/GoChat")
    subprocess.check_call('git push origin v0.1.8', cwd="/sdk/GoChat")
