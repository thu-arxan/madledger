import os
import re

ROOT = '../..'
positions = [
        ('client', 'explorer/config/cert'),
        # env
        ('client', 'env/bft/.clients/0/cert'),
        ('client', 'env/bft/.clients/1/cert'),
        ('client', 'env/bft/.clients/2/cert'),
        ('client', 'env/bft/.clients/3/cert'),
        ('client', 'env/bft/.clients/4/cert'),
        ('admin', 'env/bft/.clients/admin/cert'),
        ('peer', 'env/bft/.peers/0/cert'),
        ('peer', 'env/bft/.peers/1/cert'),
        ('peer', 'env/bft/.peers/2/cert'),
        ('peer', 'env/bft/.peers/3/cert'),
        ('orderer', 'env/bft/.orderers/0/cert'),
        ('orderer', 'env/bft/.orderers/1/cert'),
        ('orderer', 'env/bft/.orderers/2/cert'),
        ('orderer', 'env/bft/.orderers/3/cert'),
        # env/raft
        ('client', 'env/raft/.clients/0/cert'),
        ('client', 'env/raft/.clients/1/cert'),
        ('client', 'env/raft/.clients/2/cert'),
        ('client', 'env/raft/.clients/3/cert'),
        ('client', 'env/raft/.clients/4/cert'),
        ('admin', 'env/raft/.clients/admin/cert'),
        ('peer', 'env/raft/.peers/0/cert'),
        ('peer', 'env/raft/.peers/1/cert'),
        ('peer', 'env/raft/.peers/2/cert'),
        ('peer', 'env/raft/.peers/3/cert'),
        ('orderer', 'env/raft/.orderers/0/cert'),
        ('orderer', 'env/raft/.orderers/1/cert'),
        ('orderer', 'env/raft/.orderers/2/cert'),
        ('orderer', 'env/raft/.orderers/3/cert'),
        # env_local/raft
        ('client', 'env_local/raft/.clients/0/cert'),
        ('client', 'env_local/raft/.clients/1/cert'),
        ('client', 'env_local/raft/.clients/2/cert'),
        ('client', 'env_local/raft/.clients/3/cert'),
        ('client', 'env_local/raft/.clients/4/cert'),
        ('admin', 'env_local/raft/.clients/admin/cert'),
        ('peer', 'env_local/raft/.peers/0/cert'),
        ('peer', 'env_local/raft/.peers/1/cert'),
        ('peer', 'env_local/raft/.peers/2/cert'),
        ('peer', 'env_local/raft/.peers/3/cert'),
        ('orderer', 'env_local/raft/.orderers/0/cert'),
        ('orderer', 'env_local/raft/.orderers/1/cert'),
        ('orderer', 'env_local/raft/.orderers/2/cert'),
        ('orderer', 'env_local/raft/.orderers/3/cert'),
        # samples
        ('client', 'samples/clients/0/cert'),
        ('client', 'samples/clients/1/cert'),
        ('client', 'samples/clients/2/cert'),
        ('client', 'samples/clients/3/cert'),
        ('client', 'samples/clients/4/cert'),
        ('admin', 'samples/clients/admin/cert'),
        ('peer', 'samples/peers/0/cert'),
        ('peer', 'samples/peers/1/cert'),
        ('peer', 'samples/peers/2/cert'),
        ('peer', 'samples/peers/3/cert'),
        ('orderer', 'samples/orderers/0/cert'),
        ('orderer', 'samples/orderers/1/cert'),
        ('orderer', 'samples/orderers/2/cert'),
        ('orderer', 'samples/orderers/3/cert'),
    ]

if raw_input('Regenerate Root CA?') in ['y','Y']:
    os.system('python gen_cert.py ca')

for tag, pos in positions:
    path = os.path.join(ROOT, pos)
    assert(os.path.exists(path))
    os.system('python gen_cert.py auth --dir=%s --name=%s' %(path, tag))
    print(path)
