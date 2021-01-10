package com.example.gdsudao

import android.app.DatePickerDialog
import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.util.Log
import android.view.Menu
import android.view.MenuItem
import android.widget.Toast
import com.example.gdsudao.activity.MenuActivity
import kotlinx.android.synthetic.main.activity_cadastro_area.*
import kotlinx.android.synthetic.main.activity_detalhes.*
import java.text.SimpleDateFormat
import java.util.*

class DetalhesActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_detalhes)

        setSupportActionBar(findViewById(R.id.toolbarDetalhes))

        // Get the areas
        var sp = com.example.gdsudao.utils.SharedPreferences()
        val bundle = intent.extras
        var areas = sp.RecuperarListaAreas(this)
        val areaIndex = bundle?.getInt("item")

        if(areaIndex != null){
            // Toast.makeText(this, "Area: ${areaIndex}", Toast.LENGTH_SHORT).show(); FUNCIONANDO
            Log.println(Log.DEBUG, "area", "Area: ${areas[areaIndex].toString()}")

            tvDataEstimada.text = areas[areaIndex].proxcorte
            tvNomeArea.text = areas[areaIndex].nome
            tvDataUltimoPastejo.text = areas[areaIndex].dataCorte
            tvGdAcumulado.text = areas[areaIndex].st
            tvEstacao.text  = areas[areaIndex].codigoEstacao
            tvPrevisoes.text  = areas[areaIndex].previsao + "%"
            tvNormais.text  = areas[areaIndex].normal + "%"
            tvDiarios.text = areas[areaIndex].diario + "%"
            tvNumeroCortes.text = areas[areaIndex].numeroCorte

        }else{
            Toast.makeText(this, "Erro ao apresentar os detalhes.", Toast.LENGTH_SHORT).show()
        }

    }

    override fun onCreateOptionsMenu(menu: Menu?): Boolean {
        menuInflater.inflate(R.menu.detalhes_menu, menu)
        return super.onCreateOptionsMenu(menu)
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        var id = item.title

        // ADICIONAR NOVO PASTEJO
        if (id == "Adicionar pastejo"){
            // TODO: Adicionar um novo pastejo
            var cal = Calendar.getInstance()

            val dateSetListener = DatePickerDialog.OnDateSetListener { view, year, monthOfYear, dayOfMonth ->
                cal.set(Calendar.YEAR, year)
                cal.set(Calendar.MONTH, monthOfYear)
                cal.set(Calendar.DAY_OF_MONTH, dayOfMonth)

                val dataFmt = SimpleDateFormat("yyyy-MM-dd", Locale.US).format(cal.time)
                val sdf = SimpleDateFormat("dd/MM/yyyy", Locale.US)

                val sp = com.example.gdsudao.utils.SharedPreferences()
                val bundle = intent.extras
                var areas = sp.RecuperarListaAreas(this)
                val areaIndex = bundle?.getInt("item")
                if (areaIndex != null){
                    areas[areaIndex].dataCorte = dataFmt
                    val nc = areas[areaIndex].numeroCorte.toInt() + 1
                    areas[areaIndex].numeroCorte = nc.toString()

                    sp.AtualizarAreaLocal(this, areas[areaIndex], areaIndex)

                    var intent = Intent(this, MenuActivity::class.java)
                    startActivity(intent)
                }else{
                    Toast.makeText(this, "Erro ao cadastrar novo corte", Toast.LENGTH_SHORT).show()
                }
            }

            DatePickerDialog(this, dateSetListener,
                    cal.get(Calendar.YEAR),
                    cal.get(Calendar.MONTH),
                    cal.get(Calendar.DAY_OF_MONTH)).show()
        }

        // EXCLUIR AREA ATUAL
        if (id == "Excluir área"){
            val sp = com.example.gdsudao.utils.SharedPreferences()
            val bundle = intent.extras
            var areas = sp.RecuperarListaAreas(this)
            val areaIndex = bundle?.getInt("item")
            if(areaIndex != null){
                sp.RemoverAreaLista(this, areas[areaIndex])
                var intent = Intent(this, MenuActivity::class.java)
                startActivity(intent)
            }else{
                Toast.makeText(this, "Erro ao excluir area.", Toast.LENGTH_SHORT).show()
            }
        }

        return super.onOptionsItemSelected(item)
    }
}