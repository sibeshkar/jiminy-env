import getpass
import logging
import os
import random
import uuid



def default_client_id():
    return '{}-{}'.format(uuid.uuid4(), getpass.getuser())

def rewarder_session(which):
    if which is None:
        which = rewarder.RewarderSession

    if isinstance(which, type):
        return which
    else:
        raise error.Error('Invalid RewarderSession driver: {!r}'.format(which))