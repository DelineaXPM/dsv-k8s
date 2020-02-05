#!/bin/sh

program=$(basename $0)

DEFAULT_NAMESPACE="default"
DEFAULT_DIRECTORY=.
DEFAULT_BITS=2048

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

kubectl=$(which kubectl)

if test $? -ne 0
then
  echo "${program}: unable to locate kubectl executable; exiting"
  exit 1
fi

openssl=$(which openssl)

if test $? -ne 0
then
  echo "${program}: unable to locate openssl executable; exiting"
  exit 1
fi

cd "${directory}"

"${openssl}" genrsa -out "${name}.key" ${bits} 2> /dev/null

csr_conf=`mktemp -t dsv-csr.XXXXXX`
trap 'rm "${csr_conf}"; exit' SIGINT SIGTERM EXIT

cat >"${csr_conf}"<<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[v3_req]
basicConstraints = CA:FALSE
keyUsage = digitalSignature, keyEncipherment, nonRepudiation
extendedKeyUsage = serverAuth
subjectAltName = @alt_names
[alt_names]
DNS.1 = ${name}
DNS.2 = ${name}.${namespace}
DNS.3 = ${fqdn}
DNS.4 = ${fqdn}.cluster.local
EOF
"${openssl}" req -new -key "${name}.key" -subj "/CN=${fqdn}" \
	-out "${name}.csr" -config "${csr_conf}" > /dev/null

"${kubectl}" delete csr --ignore-not-found=true "${fqdn}" > /dev/null
"${kubectl}" create -f - <<EOF > /dev/null
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: ${fqdn}
spec:
  groups:
  - system:authenticated
  request: $(cat "${name}.csr" | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF

retries=0
while :; do
  "${kubectl}" get csr "${fqdn}" > /dev/null
  test "$?" -eq 0 && break
  if test $retries -gt 9
  then
    echo "no csr after ${retries} seconds; exiting"
    exit 1
  fi
  sleep 1
done

"${kubectl}" certificate approve "${fqdn}" > /dev/null

retries=0
while :; do
  cert=$("${kubectl}" get csr ${fqdn} -o jsonpath='{.status.certificate}')
  test -n "${cert}" && break
  if test $retries -gt 9
  then
    echo "no certificate after ${retries} seconds; exiting"
    exit 1
  fi
  sleep 1
done

echo ${cert} | "${openssl}" base64 -d -A -out "${name}.pem"
