// assets/js/api.js
document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('nota-form');
    const list = document.getElementById('notas-list');
    const template = document.getElementById('nota-template');
    const submitBtn = document.getElementById('nota-submit');
    
    // EXIGENCIA PEC3: Promesas y manipulación asíncrona
    const cargarNotas = () => {
        fetch('/api/notas')
            .then(response => {
                if (!response.ok) throw new Error("Error en red");
                return response.json();
            })
            .then(data => {
                list.innerHTML = '';
                if (data && data.length > 0) {
                    data.forEach(nota => renderNota(nota));
                } else {
                    list.innerHTML = '<p style="text-align:center; color:#8a7a5f;">Aún no tienes notas escritas.</p>';
                }
            })
            .catch(err => console.error("Error al cargar notas:", err));
    };

    const renderNota = (nota) => {
        // EXIGENCIA PEC3: Uso de la etiqueta <template>
        const clone = template.content.cloneNode(true);
        const article = clone.querySelector('.nota-card');
        
        clone.querySelector('.nota-titulo').textContent = nota.titulo;
        clone.querySelector('.nota-contenido').textContent = nota.contenido;
        
        // EXIGENCIA PEC3: Manipulación visual mediante clases CSS
        article.classList.add('is-loaded');
        
        clone.querySelector('.btn-editar').addEventListener('click', () => {
            document.getElementById('nota-id').value = nota.id;
            document.getElementById('nota-titulo').value = nota.titulo;
            document.getElementById('nota-contenido').value = nota.contenido;
            
            // Alteración visual por clases
            form.classList.add('is-editing');
            submitBtn.textContent = 'Actualizar Nota';
            submitBtn.style.backgroundColor = '#ffffff'; // Efecto visual leve
            
            // Scroll suave hacia arriba
            window.scrollTo({ top: form.offsetTop - 100, behavior: 'smooth' });
        });
        
        clone.querySelector('.btn-eliminar').addEventListener('click', () => {
            if(confirm('¿Seguro que deseas eliminar esta estrategia?')) {
                // DELETE con Promesa
                fetch(`/api/notas?id=${nota.id}`, { method: 'DELETE' })
                    .then(() => cargarNotas())
                    .catch(err => console.error("Error al eliminar:", err));
            }
        });
        
        list.appendChild(clone);
    };

    form.addEventListener('submit', (e) => {
        e.preventDefault();
        const id = document.getElementById('nota-id').value;
        const titulo = document.getElementById('nota-titulo').value;
        const contenido = document.getElementById('nota-contenido').value;

        const payload = { titulo, contenido };
        let url = '/api/notas';
        let method = 'POST';

        if (id) {
            payload.id = parseInt(id);
            method = 'PUT'; // Exigencia PEC3: Uso correcto de métodos HTTP
        }

        fetch(url, {
            method: method,
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        })
        .then(response => response.json())
        .then(() => {
            // Limpieza y reinicio visual
            form.reset();
            document.getElementById('nota-id').value = '';
            form.classList.remove('is-editing');
            submitBtn.textContent = 'Guardar Nota';
            submitBtn.style.backgroundColor = '#c9ad7a';
            cargarNotas();
        })
        .catch(err => console.error("Error al guardar:", err));
    });

    // Carga inicial
    cargarNotas();
});