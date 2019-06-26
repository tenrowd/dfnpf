#!/usr/bin/env python3

from noise.connection import NoiseConnection, Keypair
import logging
import sys
import os
import binascii

params = dict (staticInitKey = "", ephemInitKey = "", remoteInitStatKey = "", 
        staticRespKey = "", ephemRespKey = "", remoteRespStatKey = "",  
        prologue = "", presharedKey = "")

def _set_args(connection):
    role = 'Init' if connection.noise_protocol.initiator else 'Resp'
    setters = [
        (connection.set_keypair_from_private_bytes, Keypair.STATIC, "static"+role+"Key"),
        (connection.set_keypair_from_private_bytes, Keypair.EPHEMERAL, "ephem"+role+"Key"),
        (connection.set_keypair_from_public_bytes, Keypair.REMOTE_STATIC, "remote"+role+"StatKey")       
    ]
    for func, keypair, name in setters:
        if params[name] != "":
            func(keypair, params[name])
    if params['prologue'] != "":
        connection.set_prologue(params['prologue'])
    if params['presharedKey'] != "":
        connection.set_psks(params['presharedKey'])
    return connection


def set_handsh(args):
    payload = binascii.unhexlify(args.pop())    

    initiator = NoiseConnection.from_name(args[1].encode())
    responder = NoiseConnection.from_name(args[1].encode())
    initiator.set_as_initiator()
    responder.set_as_responder()

    handshakeparams = args[2].split('_')
    i = 3
    for key in handshakeparams:
        params [key] = bytes.fromhex(args[i])
        i += 1

    initiator = _set_args (initiator)
    responder = _set_args (responder)

    return initiator, responder, payload