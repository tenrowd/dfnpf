#!/usr/bin/env python3

from noise.connection import NoiseConnection, Keypair
import logging
import sys
import os
import binascii

from iniths.init import *

import warnings


warnings.filterwarnings("ignore", message='One of ephemeral keypairs')

if __name__ == '__main__':
    if len(sys.argv) < 4 or len(sys.argv) > 12:
        print("Please provide Noise Protocol name, initiator's static and ephemeral keys,",
		    	"\n responder's static and ephemeral keys, remote key, prologue, preshared key, message",
                    "\nthere are %s of args",len(sys.argv))
        sys.exit(1)


    initiator , responder, payload = set_handsh(sys.argv)

    initiator.start_handshake()
    responder.start_handshake()

    #message = initiator. write_message()
    #recieved = responder.read_message(message)
    #message = responder.write_message()
    #recieved = initiator.read_message(message)


    ciphertext = initiator.write_message(payload)
    plaintext = responder.read_message(ciphertext)
    print(ciphertext.hex())