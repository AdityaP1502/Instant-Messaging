import os
import os.path
import sys

from process_handler import ProcessHandler


if __name__ == "__main__":
    path = os.path.join(os.path.dirname(os.path.abspath(sys.argv[0])), "doodle/spawn_process.py")
    cmd = "python {}".format(path)
    p_handler = ProcessHandler(cmd=cmd)
    p_handler.run()
    p_handler.join()
    