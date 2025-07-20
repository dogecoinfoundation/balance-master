 
 

create table if not exists utxos (
    tx_id TEXT NOT NULL,
    vout INTEGER NOT NULL,
    address TEXT NOT NULL,
    amount BIGINT NOT NULL,
    PRIMARY KEY (tx_id, vout)
);