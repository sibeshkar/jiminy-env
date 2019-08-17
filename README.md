Jiminy internal communication protocols
*****************************************

Use Plugin based system
====================

Here's the instructions to run this : 

1. Git clone `github.com/sibeshkar/jiminy-env`
2. Run : `make docker-run-d` to run detached mode
3. Install jiminy from `github.com/sibeshkar/jiminy`, mods branch
5. Outside container, run: `python misc/wob_new_dom.py`

To debug inside container:

1. Git clone `github.com/sibeshkar/jiminy-env`
2. Run : `make docker-run` to run interactive mode
3. In the container run : `jiminy install wob-v0.zip`
4. After installing plugin, run:  `jiminy run sibeshkar/wob-v0`
5. Outside container, run: `python misc/wob_new_dom.py`

Compiling:
1. `make docker` to compile env base, plugin-1 and plugin-2, and make docker image $VERSION
2. `make plugin-v1` to plugin-1
3. `make plugin-v2` to plugin-2

Recording:
1. `make docker-record` to enter docker container. Install plugin using `jiminy install <plugin zip file>`
2. Run `jiminy record <repo/env/task>` for e.g. `jiminy record sibeshkar/wob-v1/ClickButton`
3. Connect to `0.0.0.0:5901` on your computer with VNCviewer. Password : `boxware`
4. To use Recordings, copy the recordings from the docker container using the following command:
`docker cp <container_name>:/root/.jiminy/plugins/sibeshkar/wob-v1/recordings/* .`
5. Process using `jiminy.demonstration.event_readers.VNCDemontration` object to get iterator.

See Makefile for more commands.

Update protobuf definition :

```sh
$ protoc -I proto/ proto/env.proto --go_out=plugins=grpc:proto/
```

For Python:

```sh
$ python -m grpc_tools.protoc -I ./proto/ --python_out=./plugin-python/ --grpc_python_out=./plugin-python/ ./proto/env.proto
```


Network architecture
====================

A Jiminy environment consists of two components that run in
separate processes and communicate over the network.  The agent's
machine runs the environment **client** code, which connects
to the **remote** environment server.

Each remote exposes two ports:

- A VNC port (5900 by default). The remote runs an off-the-shelf VNC
  server (usually TigerVNC), so that users can connect their own
  VNC viewers to the environments for interactive use. VNC delivers
  pixel observations to the agent, and the agent submits keyboard and
  mouse inputs over VNC as well.

- A rewarder port (15900 by default). The rewarder protocol is a
  bi-directional JSON protocol runs over WebSockets. The rewarder
  channel provides more than just a reward signal; in addition, it allows the
  agent to submit control commands (such as to indicating which of
  the available environments should be active for a given runtime) and
  to receive structured information from the environment (such as latencies
  and performance timings).

VNC system and Remote Frame Buffer protocol
===========================================
 
Keyboard and mouse actions and pixel observations are sent between the
client and the remote using the `VNC
<https://en.wikipedia.org/wiki/Virtual_Network_Computing>`__
system. VNC is a pervasive standard for remote desktop operation. Many
implementations of VNC are available online, including VNC viewers
that make it easy to observe a running agent.

More information about the Remote Frame Buffer protocol can be found
in the official `IETF RFC <https://tools.ietf.org/html/rfc6143>`__
spec, and in other tutorials elsewhere on the internet.

Rewarder protocol
=================

The Rewarder protocol is a Jiminy-specific bi-directional JSON
protocol run over WebSockets. In addition to the actions and
observations provided by the VNC connection, the rewarder connection
allows the agent to submit control commands to the environment, and to
receive rewards and other informational messages. This section details
the format of the Rewarder protocol.

Message format
--------------

Each message is serialized as a JSON object with the following
structure:

.. code::
		  
    {
      "method": [string],
      "headers": [object],
      "body": [object]
    }

For example, a ``v0.env.describe`` message might look as follows:

.. code::

    {
      "method": "v0.env.describe",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
        "episode_id": "1.2",
      },
      "body": {
        "env_id": "internet.SlitherIO-v0",
        "env_state": "running",
        "fps": 60
      }
    }


Each message should have a unique ``message_id`` header and a ``sent_at``
header (which should be the current UNIX timestamp with at least
millisecond precision).

Server to Client messages
-------------------------

env.describe
~~~~~~~~~~~~

.. code:: 
		  
    {
      "method": "v0.env.describe",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
        "episode_id": "1.2",
      },
      "body": {
        "env_id": "internet.SlitherIO-v0",
        "env_state": "running",
      }
    }

env.reward
~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.env.reward",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
        "episode_id": "1.2",
      },
      "body": {
        "reward": -3,
        "done": False,
    	"info": {},
      }
    }



env.text
~~~~~~~~

.. code::
		  
    {
      "method": "v0.env.text",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
        "episode_id": "1.2",
      },
      "body": {
        "text": "this is some text"
      }
    }

env.observation
~~~~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.env.observation",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
        "episode_id": "1.2"
      },
      "body": {
        "observation": [0.12, 0.51, 2, 12]
      }
    }

connection.close
~~~~~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.connection.close",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15
      },
      "body": {
        "message": "Disconnected since time limit reached"
      }
    }

reply.error
~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.reply.error",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
    	"parent_message_id": "26"
      },
      "body": {
        "message": "No such environment: llama"
      }
    }

reply.env.reset
~~~~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.reply.env.reset",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
    	"parent_message_id": "26",
    	"episode_id": "1.2"
    	
      },
      "body": {}
    }

reply.env.launch
~~~~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.reply.env.launch",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
    	"parent_message_id": "26",
      },
      "body": {}
    }
    
reply.control.ping
~~~~~~~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.reply.control.ping",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15,
    	"parent_message_id": "26"
      },
      "body": {}
    }

Client to server messages
-------------------------

agent.action
~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.agent.action",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15
      },
      "body": {
        "action: [["JoystickAxisXEvent", 0.1],
                  ["JoystickAxisZEvent", 0.1]]
      }
    }

env.launch
~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.env.launch",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15
      },
      "body": {
        "env_id": "sibeshkar/wob-v0"
        "fps" : "16"
      }
    }


env.reset
~~~~~~~~~

.. code::
		  
    {
      "method": "v0.env.reset",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15
      },
      "body": {
        "env_id": "sibeshkar/wob-v0",
        "task_id": "TicTacToe"
      }
    }


control.ping
~~~~~~~~~~~~

.. code::
		  
    {
      "method": "v0.control.ping",
      "headers": {
        "sent_at": 1479493678.1937322617,
        "message_id": 15
      },
      "body": {}
    }
