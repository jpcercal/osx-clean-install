#!/usr/bin/env python

import argparse
import asyncio
import os
import time
import sys

parser = argparse.ArgumentParser()
parser.add_argument("-H", "--host", type=str, default='127.0.0.1', help='Hostname', required=True)
parser.add_argument("-P", "--port", type=int, default=80, help='Port', required=True)
parser.add_argument("-T", "--timeout", type=int, default=1, help='Timeout', required=False)
parser.add_argument("-D", "--delay", type=int, default=2, help='Delay', required=False)
args = parser.parse_args()

async def wait_for_port(host, port, timeout, delay):
    tmax = time.time() + timeout
    while time.time() < tmax:
        try:
            _reader, writer = await asyncio.wait_for(asyncio.open_connection(host, port), timeout=5)
            writer.close()
            await writer.wait_closed()
            return True
        except:
            if delay:
                await asyncio.sleep(delay)
    return False

is_port_running = asyncio.run(wait_for_port(
    args.host, 
    args.port,
    int(args.timeout),
    int(args.delay),
))

sys.exit(os.EX_OK) if is_port_running else sys.exit(os.EX_TEMPFAIL)
