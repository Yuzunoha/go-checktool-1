<?php

function p($a)
{
    echo $a . PHP_EOL;
}

for ($i = 1; $i <= 100; $i++) {
    // 場合には、数の代わりに「3の倍数であり、5の倍数
    if ($i % 15 === 0) {
        p("3の倍数であり、5の倍数");
    } else if ($i % 3 === 0) {
        p("3の倍数");
    } else if ($i % 5 === 0) {
        p("5の倍数");
    } else {
        p($i);
    }
}
