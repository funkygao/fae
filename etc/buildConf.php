<?php

if ($argc != 7) {
    die("Usage: $argv[0] [etcd_servers] [http_listen_addr] [pprof_listen_addr] [metrics_logfile] [rpc.listen_addr] [servants.idgen_worker_id]\n");
}

$target = 'faed.cf.rc';
$template = 'faed.cf.sample';
echo "reading $template\n";
$body = file_get_contents($template);
$search = array(
    '{etcd_servers}',
    '{http_listen_addr}',
    '{pprof_listen_addr}',
    '{metrics_logfile}',
    '{rpc_listen_addr}',
    '{idgen_worker_id}'
);

$server = str_replace(',', "\",\n       \"", $argv[1]);
$replace = array(
    '"'.$server.'"',
    '"'.$argv[2].'"',
    '"'.$argv[3].'"',
    '"'.$argv[4].'"',
    '"'.$argv[5].'"',
    $argv[6]
);

$body = str_replace($search, $replace, $body);
file_put_contents($target, $body);
echo "write to $target\n";
