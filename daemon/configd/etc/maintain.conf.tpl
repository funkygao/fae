<?php

return array(
    "maintain_mode" => array(
        {{range $index, $value := .}}"{{$index}}" => {{$value}},
        {{end}}
    ),
);
