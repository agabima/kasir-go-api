const API_URL = "http://localhost:8080/produk";

const produkForm = document.getElementById("produkForm");
const produkTable = document.getElementById("produkTable");

// 🔹 Load data saat halaman dibuka
document.addEventListener("DOMContentLoaded", loadProduk);

async function loadProduk() {
  produkTable.innerHTML = `<tr><td colspan="6" class="text-center p-3">Loading...</td></tr>`;
  
  const res = await fetch(API_URL);
  const data = await res.json();
    //alert(data.length);
  if (data.length === 0) {
    produkTable.innerHTML = `<tr><td colspan="6" class="text-center p-3">Belum ada produk.</td></tr>`;
    return;
  }

  produkTable.innerHTML = "";
  data.forEach((p) => {
    produkTable.innerHTML += `
      <tr class="border-b hover:bg-gray-50">
        <td class="py-2 px-3 border">${p.id}</td>
        <td class="py-2 px-3 border">${p.nama}</td>
        <td class="py-2 px-3 border">Rp${p.harga.toLocaleString()}</td>
        <td class="py-2 px-3 border">${p.stok}</td>
        <td class="py-2 px-3 border">${p.kategori_id}</td>
        <td class="py-2 px-3 border text-center">
          <button onclick="hapusProduk(${p.id})" class="text-red-600 hover:underline">Hapus</button>
        </td>
      </tr>`;
  });
}

// 🔹 Tambah produk baru
produkForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const newProduk = {
    nama: document.getElementById("nama").value,
    harga: parseFloat(document.getElementById("harga").value),
    stok: parseInt(document.getElementById("stok").value),
    kategori_id: parseInt(document.getElementById("kategori_id").value)
  };

  const res = await fetch(API_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(newProduk),
  });

  if (res.ok) {
    alert("Produk berhasil ditambahkan!");
    produkForm.reset();
    loadProduk();
  } else {
    const err = await res.json();
    alert("Gagal menambah produk: " + err.error);
  }
});

// 🔹 Hapus produk
async function hapusProduk(id) {
  if (!confirm("Yakin ingin menghapus produk ini?")) return;

  const res = await fetch(`${API_URL}/${id}`, { method: "DELETE" });
  if (res.ok) {
    alert("Produk berhasil dihapus!");
    loadProduk();
  } else {
    alert("Gagal menghapus produk.");
  }
}
