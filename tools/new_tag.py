import subprocess

if __name__ == '__main__':
    subprocess.check_call('git tag -a v0.1.6 -m "Release_v0.1.6"')
    subprocess.check_call('git push origin v0.1.6')