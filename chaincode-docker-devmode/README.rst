Using dev mode
==============

> **Copyright before 06/08/2018 Hyperledger Fabric** [Link to source code](https://github.com/hyperledger/fabric-samples/tree/release-1.1/chaincode-docker-devmode) 
> Licensed under Creative Commons Attribution 4.0 International License 
> **Copyright after 06/08/2018 Owkin, Substra project**, licensed under the Apache-2.0. 

Normally chaincodes are started and maintained by peer. However in “dev
mode", chaincode is built and started by the user. This mode is useful
during chaincode development phase for rapid code/build/run/debug cycle
turnaround.

We start "dev mode" by leveraging pre-generated orderer and channel artifacts for
a sample dev network.  As such, the user can immediately jump into the process
of compiling chaincode and driving calls.

Install Fabric Samples
----------------------

If you haven't already done so, please install the doc [samples](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html).

Navigate to the ``chaincode-docker-devmode`` directory of the ``fabric-samples``
clone:

.. code:: bash

  cd chaincode-docker-devmode

Download docker images
^^^^^^^^^^^^^^^^^^^^^^

We need four docker images in order for "dev mode" to run against the supplied
docker compose script.  If you installed the ``fabric-samples`` repo clone and
followed the instructions to [download-platform-specific-binaries](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#download-platform-specific-binaries), then
you should have the necessary Docker images installed locally.

.. note:: If you choose to manually pull the images then you must retag them as
          ``latest``.

Issue a ``docker images`` command to reveal your local Docker Registry.  You
should see something similar to following:

.. code:: bash

  docker images
  REPOSITORY                     TAG                                  IMAGE ID            CREATED             SIZE
  hyperledger/fabric-tools       latest                c584c20ac82b        9 days ago         1.42 GB
  hyperledger/fabric-tools       x86_64-1.1.0-preview  c584c20ac82b        9 days ago         1.42 GB
  hyperledger/fabric-orderer     latest                2fccc91736df        9 days ago         159 MB
  hyperledger/fabric-orderer     x86_64-1.1.0-preview  2fccc91736df        9 dyas ago         159 MB
  hyperledger/fabric-peer        latest                337f3d90b452        9 days ago         165 MB
  hyperledger/fabric-peer        x86_64-1.1.0-preview  337f3d90b452        9 days ago         165 MB
  hyperledger/fabric-ccenv       latest                82489d1c11e8        9 days ago         1.35 GB
  hyperledger/fabric-ccenv       x86_64-1.1.0-preview  82489d1c11e8        9 days ago         1.35 GB

.. note:: If you retrieved the images through the [download-platform-specific-binaries](http://hyperledger-fabric.readthedocs.io/en/latest/samples.html#download-platform-specific-binaries),
          then you will see additional images listed.  However, we are only concerned with
          these four.

Now open three terminals and navigate to your ``chaincode-docker-devmode``
directory in each.

Terminal 1 - Start the network
------------------------------

.. code:: bash

    docker-compose -f docker-compose-simple.yaml up

The above starts the network with the ``SingleSampleMSPSolo`` orderer profile and
launches the peer in "dev mode".  It also launches two additional containers -
one for the chaincode environment and a CLI to interact with the chaincode.  The
commands for create and join channel are embedded in the CLI container, so we
can jump immediately to the chaincode calls.

Terminal 2 - Build & start the chaincode
----------------------------------------

.. code:: bash

  docker exec -it chaincode bash

You should see the following:

.. code:: bash

  root@d2629980e76b:/opt/gopath/src/chaincode#

Now, compile your chaincode:

.. code:: bash

  go build -o substra_chaincode

Now run the chaincode:

.. code:: bash

  CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=mycc:0 ./substra_chaincode

The chaincode is started with peer and chaincode logs indicating successful registration with the peer.
Note that at this stage the chaincode is not associated with any channel. This is done in subsequent steps
using the ``instantiate`` command.

Terminal 3 - Use the chaincode
------------------------------

Even though you are in ``--peer-chaincodedev`` mode, you still have to install the
chaincode so the life-cycle system chaincode can go through its checks normally.
This requirement may be removed in future when in ``--peer-chaincodedev`` mode.

We'll leverage the CLI container to drive these calls.

.. code:: bash

  docker exec -it cli bash

.. code:: bash

  peer chaincode install -p chaincodedev/chaincode/ -n mycc -v 0
  peer chaincode instantiate -n mycc -v 0 -c '{"Args":["init"]}' -C myc

Now issue an invoke to create a problem.

.. code:: bash

  peer chaincode invoke -n mycc -c '{"Args":["addProblem","5c1d9cd1c2c1082dde0921b56d11030c81f62fbb51932758b58ac2569dd0b379", "msi classification","https://toto/problem/222/description","https://toto/problem/222/metrics", "fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8d","data_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8a, data_fd1bb7c31f62244c0f3a761cc168804227115793d01c270021fe3f7935482d8e","all"]}' -C myc

Finally, query all problems.  We should see the newly added problem.

.. code:: bash

  peer chaincode query -n mycc -c '{"Args":["queryProblem"]}' -C myc


