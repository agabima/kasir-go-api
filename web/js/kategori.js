const api = "http://localhost:8080/kategori";

async function loadKategori() {
  const res = await fetch(api);
  const data = await res.json();
  const tbody = document.getElementById("kategoriTable");
  tbody.innerHTML = data
    .map(
      (k) => `
      <tr>
        <td class="border p-2">${k.id}</td>
        <td class="border p-2">${k.nama}</td>
        <td class="border p-2">
          <button onclick="hapusKategori(${k.id})" class="text-red-500 hover:underline">Hapus</button>
        </td>
      </tr>
    `
    )
    .join("");
}

async function tambahKategori() {
  const nama = document.getElementById("namaKategori").value;
  if (!nama) return alert("Nama kategori tidak boleh kosong!");

  await fetch(api, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ nama }),
  });
  document.getElementById("namaKategori").value = "";
  loadKategori();
}

async function hapusKategori(id) {
  if (!confirm("Yakin ingin menghapus?")) return;
  await fetch(`${api}/${id}`, { method: "DELETE" });
  loadKategori();
}

loadKategori();
