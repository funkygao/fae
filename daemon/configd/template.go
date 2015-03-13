package main

const (
	FAE_TPL = `<?php

return array(
    "config" => array(
        "send_timeout" => 4000, // ms
        "recv_timeout" => 4000, // ms
        "write_buffer" => 1024, // byte
        "read_buffer" => 1024,  // byte
        "retries" => 1,
    ),

    "hosts" => array(
    	{{range $index, $value := .Servers}}"{{$value}}",
    	{{end}}
    ),

    "ports" => array(
        {{range $index, $value := .Ports}}{{$value}},
    	{{end}}
    ),
);`

	MAINTAIN_TPL = `<?php

return array(
    "maintain_mode" => array(
        {{range $index, $value := .}}"{{$index}}" => {{$value}},
        {{end}}
    ),
);`
)
