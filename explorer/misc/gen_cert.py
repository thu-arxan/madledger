#!/usr/bin/python
import os
import argparse

def check_dir_exist(d):
    if not os.path.exists(d):
        print('mkdir %s' % d)
        os.mkdir(d)

def check_file_not_exist(f):
    print('check', f)
    if os.path.exists(f) or os.path.islink(f):
        print('rm %s' % f)
        os.remove(f)

def args_to_dict(args):
    res = {}
    argsDict = args.__dict__;
    for eachArg in argsDict.keys():
        if(type(argsDict[eachArg]) != "<type 'string'>"):
            res[eachArg] = argsDict[eachArg]
        else:
            pass
    return res

def execute(cmd):
    print('\033[32m%s\033[0m' % cmd)
    if os.system(cmd):
        print('command failed.')

def relative_path(a, b): 
    a = os.path.abspath(a)
    b = os.path.abspath(b)
    a, b = a.split('/'), b.split('/')
    intersection = 0
    for index in range(min(len(a), len(b))):
        m, n = a[index], b[index]
        if m != n:
            intersection = index
            break
    def backward():
        return (len(a) - intersection) * '../'
    
    def forward():
        return '/'.join(b[intersection:])
    
    out = backward() + forward()
    return out

def generate_ca(args):
    args = args_to_dict(args)
    args['ca_keyfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.key')
    args['ca_pemfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.pem')

    check_dir_exist(args['ca_dir'])
    check_file_not_exist(args['ca_keyfile'])
    check_file_not_exist(args['ca_pemfile'])

    if args['method'] == 'ecc':
        execute('openssl ecparam -out {ca_keyfile} -name prime256v1 -genkey'.format(**args))
    elif args['method'] == 'dsa':
        execute('openssl genrsa -des3 -passout pass:"{dsa_password}" -out {ca_keyfile} 2048'.format(**args))
    execute('openssl req -x509 -new -nodes -key {ca_keyfile} -passin pass:"{dsa_password}" -config {ca_config} -sha256 -days 500 -subj "{ca_subject}" -out {ca_pemfile}'.format(**args))


def auth_client(args):
    args = args_to_dict(args)
    args['ca_keyfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.key')
    args['ca_pemfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.pem')
    args['ca_pemfile_abs'] = os.path.abspath(args['ca_pemfile'])
    args['ca_pemfile_rel'] = './'+os.path.join(relative_path(args['dir'], args['ca_dir']), args['ca_filename']+'.pem')
    args['keyfile'] = os.path.join(args['dir'], args['name'] + '.key')
    args['crtfile'] = os.path.join(args['dir'], args['name'] + '.crt')
    args['csrfile'] = os.path.join(args['dir'], args['name'] + '.csr')
    args['ca_pemlink'] = os.path.join(args['dir'], args['ca_link'])
    print(args['ca_pemfile_rel'])

    assert(os.path.exists(args['dir']))
    check_file_not_exist(args['keyfile'])
    check_file_not_exist(args['crtfile'])
    check_file_not_exist(args['csrfile'])
    check_file_not_exist(args['ca_pemlink'])
    check_file_not_exist(os.path.join(args['dir'], 'ca.cer'))
    check_file_not_exist(os.path.join(args['dir'], args['name']+'.pem'))
    execute('cp {ca_pemfile_abs} {ca_pemlink}'.format(**args))

    if args['method'] == 'ecc':
        execute('openssl ecparam -out {keyfile} -name prime256v1 -genkey'.format(**args))
    elif args['method'] == 'dsa':
        execute('openssl genrsa -out {keyfile} 2048'.format(**args))

    execute('openssl req -new -key {keyfile} -out {csrfile} -config {ca_config} -subj "{ca_subject}"'.format(**args))

    execute('openssl x509 -req -in {csrfile} -CA {ca_pemfile} -CAkey {ca_keyfile} -CAcreateserial -out {crtfile} -days 400 -sha256 -extfile {config} -passin pass:"{dsa_password}"'.format(**args))

    execute('cat {ca_pemfile} >> {crtfile}'.format(**args))

    pass


def list_ca(args):
    args = args_to_dict(args)
    args['ca_keyfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.key')
    args['ca_pemfile'] = os.path.join(args['ca_dir'], args['ca_filename'] + '.pem')
    
    execute('openssl x509 -text -noout -in {ca_pemfile}'.format(**args))

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--method', help='ecryption method', choices=['dsa','ecc'], default='ecc')
    parser.add_argument('--ca_dir', help='directory where CA store.', default='./CA')
    parser.add_argument('--ca_filename', help='base name of CA files.', default='localhostCA')
    parser.add_argument('--dsa_password', type=str, help='dsa_password of your CA\' private key.', default='unsafe dsa_password')
    parser.add_argument('--ca_config', help='path to CA configures.', default='./config/localhostCA.conf')
    parser.add_argument('--ca_subject', help='Extensions of CA generations', default='/C=CN/ST=Beijing/L=./O=MadLedger')
    subparsers = parser.add_subparsers(help='sub-command help')

    parser_ca = subparsers.add_parser('ca', help='create the top-level CA')
    parser_ca.set_defaults(func=generate_ca)

    parser_list = subparsers.add_parser('list', help='list CA infos')
    parser_list.set_defaults(func=list_ca)

    parser_auth = subparsers.add_parser('auth', help='Generate client CA cert.')
    parser_auth.add_argument('--dir', help='Directory where cert store.', default='./cert')
    parser_auth.add_argument('--name', help='Base name of Cert files.', default='client')
    parser_auth.add_argument('--ca_link', help='Soft link of CA\'s pem file.', default='CA.pem')
    parser_auth.add_argument('--config', help='Config file of the cert generation.', default='./config/localhost.conf')
    parser_auth.set_defaults(func=auth_client)

    args = parser.parse_args()
    print(args)
    args.func(args)
