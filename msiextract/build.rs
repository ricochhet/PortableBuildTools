extern crate embed_resource;
fn main() {
    embed_resource::compile("msiextract.rc", embed_resource::NONE);
}
