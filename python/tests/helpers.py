import sys
import os

THIS_FILE_DIR = os.path.dirname(os.path.realpath(__file__))
context_dir = os.path.abspath(os.path.join(THIS_FILE_DIR, ".."))
sys.path.append(context_dir)

import imp
context = imp.load_source('context', 'pycontext')
