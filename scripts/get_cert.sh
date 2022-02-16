#!/bin/sh

program=$(basename $0)

DEFAULT_NAMESPACE="default"
DEFAULT_DIRECTORY=.
DEFAULT_BITS=4096

usage() {
  cat<<EOF
Usage: ${program} -n NAME [OPTIONS]...

        -n, -name, --name NAME
                Maps to the host portion of the FQDN that is the subject of the
                certificate; also the basename of the certificate and key files.
        -d, -directory, --directory=DIRECTORY
                The location of the resulting certificate and private-key. The
                default is '${DEFAULT_DIRECTORY}'
        -N, -namespace, --namespace=NAMESPACE
                Represents the Kubernetes cluster Namespace and maps to the
                domain of the FQDN that is the subject of the certificate.
                the default is '${DEFAULT_NAMESPACE}'
        -b, -bits, --bits=BITS
                the RSA key size in bits; default is ${DEFAULT_BITS}
EOF
}

options=$(getopt -l "bits:,directory:,name:,namespace:" -o "b:d:n:N:" -a -n "${program}" -- "$@")

test $? -ne 0 && usage && exit 1

eval set -- ${options}

while true
do
  case "${1}" in
  --name|-n)
    shift
    name="${1}"
    ;;
  --namespace|-N)
    shift
    namespace="${1}"
    ;;
  --bits|-b)
    shift
    bits="${1}"
    ;;
  --directory|-d)
    shift
    directory="${1}"
    ;;
  --)
    shift
    break
    ;;
  *)
    usage
    exit 1
    ;;
  esac
  shift
done

if test -z "${name}"
then
  usage
  exit 1
fi

test -z "${bits}" && bits=${DEFAULT_BITS}
test -z "${directory}" && directory=${DEFAULT_DIRECTORY}
test -z "${namespace}" && namespace=${DEFAULT_NAMESPACE}

fqdn="${name}.${namespace}.svc"

openssl=$(which openssl)

if test $? -ne 0
then
  echo "${program}: unable to locate openssl executable; exiting"
  exit 1
fi

cd "${directory}"

req_conf=`mktemp -t dsv-req.conf.XXXXXX`
trap 'rm "${req_conf}"; exit' SIGINT SIGTERM EXIT

cat >"${req_conf}"<<EOF
[req]
default_bits       = ${DEFAULT_BITS}
distinguished_name = dn
prompt             = no

[dn]
organizationName   = system:nodes
commonName         = system:node:${fqdn}

[ext]
subjectAltName     = @sans

[sans]
DNS.1              = ${name}
DNS.2              = ${fqdn}
EOF

"${openssl}" req -new -x509 -sha256 -nodes -keyout ${name}.key -days $((365*5))\
  -extensions ext -out ${name}.pem -config $req_conf 2>/dev/null
