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
 
    ciphertext = initiator.write_message(payload)
    ckI = initiator.noise_protocol.handshake_state.symmetric_state.ck
    
    plaintext = responder.read_message(ciphertext)
    ckR = responder.noise_protocol.symmetric_state.ck

    print(ckI.hex())
    print(ckR.hex())